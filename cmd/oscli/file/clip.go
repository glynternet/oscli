package file

import (
	"encoding/json"
	"fmt"
	"github.com/glynternet/go-osc/osc"
	"io"
	"log"
	"time"

	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCombinedClipsFile = "./combined-clips.osc"

func Clip(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		output string

		cmd = &cobra.Command{
			Use:   "clip",
			Short: "not to be commited to mainline oscli",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, oscFiles []string) error {
				logger.Print("starting")
				var nes []namedEntries
				for _, f := range oscFiles {
					logger.Printf("About to read from %s\n", f)
					r, err := readFromFile(logger, f)
					if err != nil {
						return errors.Wrapf(err, "reading recording from file:%s", f)
					}
					logger.Printf("Recording read from %s\n", f)
					nes = append(nes, namedEntries{
						name:    f,
						Entries: r.Data.Entries,
					})
				}
				clips, err := clipAllEntries(logger, nes...)
				if err != nil {
					return errors.Wrap(err, "clipping all entries")
				}
				combined, err := combineClips(logger, clips...)
				if err != nil {
					return errors.Wrap(err, "combining all clips")
				}

				wc, err := fileCreatingWriteCloser(logger, output)
				if err != nil {
					return errors.Wrap(err, "creating WriteCloser")
				}
				return catchFirstLogOthers(logger, writeToWriteCloser(combined, wc)...)
			},
		}
	)

	parent.AddCommand(cmd)
	cmd.Flags().StringVar(&output, "output-file", defaultCombinedClipsFile, "file to write combined osc content to")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}

func clipAllEntries(logger *log.Logger, nes ...namedEntries) ([]namedEntries, error) {
	var clips []namedEntries
	for _, r := range nes {
		logger.Printf("Clipping %s", r.name)
		clipped, err := clipEntries(logger, r.Entries)
		if err != nil {
			return nil, errors.Wrapf(err, "clipping entries for recording:%s", r.name)
		}
		clips = append(clips, namedEntries{
			name:    r.name,
			Entries: clipped,
		})
	}
	return clips, nil
}

func combineClips(logger *log.Logger, nes ...namedEntries) (combinedClips, error) {
	channelMetadata := make(channelMetadata)
	var combined record.Entries
	for i, ne := range nes {
		logger.Printf("Combining %s", ne.name)
		combined = record.Combine(combined, fixLoudnessChannel(ne.Entries, int32(i)))
		channelMetadata[i] = ne.name
	}
	return combinedClips{
		ChannelMetadata: channelMetadata,
		Recording:       record.Recording{Data: record.RecordingData{Entries: combined}},
	}, nil
}

func fixLoudnessChannel(entries record.Entries, channel int32) record.Entries {
	var fixed record.Entries
	entries.ForEach(func(i int, entry record.Entry) {
		msg := entry.Packet.(*osc.Message)
		if !isLoudnessMessage(msg) {
			panic("not loudness message")
		}
		fixedMsg := msg
		fixedMsg.Arguments[1] = channel
		fixed = append(fixed, record.Entry{
			Elapsed: entry.Elapsed,
			Packet:  fixedMsg,
		})
	})
	return fixed
}

type namedEntries struct {
	name string
	record.Entries
}

type combinedClips struct {
	ChannelMetadata channelMetadata  `json:"channel_metadata"`
	Recording       record.Recording `json:"recording"`
}

func (c combinedClips) WriteTo(w io.Writer) (n int64, err error) {
	return 0, json.NewEncoder(w).Encode(c)
}

type channelMetadata map[int]string

func clipEntries(logger *log.Logger, entries record.Entries) (record.Entries, error) {
	var playStart time.Duration
	var started bool
	var clipped record.Entries
	for _, entry := range entries {
		msg, ok := entry.Packet.(*osc.Message)
		if !ok {
			return nil, fmt.Errorf("expected *osc.Message but got %T", entry.Packet)
		}
		var hasStopped bool
		switch {
		case isPlayStartMsg(msg):
			playStart = entry.Elapsed
			logger.Printf("play start: %s", playStart)
			started = true
		case isPlayStoppedMsg(msg):
			logger.Printf("play stopped: %s, duration: %s", entry.Elapsed, entry.Elapsed-playStart)
			hasStopped = true
		case isLoudnessMessage(msg):
			if started {
				clipped = append(clipped, record.Entry{
					Elapsed: entry.Elapsed - playStart,
					Packet:  msg,
				})
			}
		default:
			return nil, fmt.Errorf("expected file_playing, file_stopped or loudness message but receieved:%+v", msg)
		}
		if hasStopped {
			break
		}
	}
	return clipped, nil
}

func isPlayStartMsg(msg *osc.Message) bool {
	return isAudioMessage(msg) && len(msg.Arguments) == 1 && msg.Arguments[0].(string) == "file_playing"
}

func isPlayStoppedMsg(msg *osc.Message) bool {
	return isAudioMessage(msg) && len(msg.Arguments) == 1 && msg.Arguments[0].(string) == "file_stopped"
}

func isLoudnessMessage(msg *osc.Message) bool {
	return isAudioMessage(msg) && len(msg.Arguments) == 3 && msg.Arguments[0].(string) == "loudness"
}

func isAudioMessage(msg *osc.Message) bool {
	return msg.Address == "/audio"
}
