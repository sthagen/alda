package cmd

import (
	log "alda.io/client/logging"
	"alda.io/client/system"
	"alda.io/client/transmitter"
	"github.com/spf13/cobra"
)

func init() {
	shutdownCmd.Flags().StringVarP(
		&playerID, "player-id", "i", "", "The ID of the player process to shut down",
	)

	shutdownCmd.Flags().IntVarP(
		&playerPort, "port", "p", -1, "The port of the player process to shut down",
	)
}

var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shut down background processes",
	RunE: func(_ *cobra.Command, args []string) error {
		players := []system.PlayerState{}

		// Determine the players to which to send a "shutdown" message based on the
		// provided CLI options.
		switch {
		// Port is explicitly specified, so use that port.
		case playerPort != -1:
			players = append(players, system.PlayerState{
				ID: "unknown", State: "unknown", Port: playerPort,
			})
		// Player ID is specified; look up the player by ID and use its port.
		case playerID != "":
			player, err := system.FindPlayerByID(playerID)
			if err != nil {
				return err
			}
			players = append(players, player)
		// Send a "shutdown" message to all known players.
		default:
			knownPlayers, err := system.ReadPlayerStates()
			if err != nil {
				return err
			}
			players = append(players, knownPlayers...)
		}

		for _, player := range players {
			transmitter := transmitter.OSCTransmitter{Port: player.Port}
			if err := transmitter.TransmitShutdownMessage(0); err != nil {
				return err
			}

			log.Info().
				Interface("player", player).
				Msg("Sent \"shutdown\" message to player process.")
		}

		return nil
	},
}
