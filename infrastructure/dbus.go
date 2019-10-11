package infrastructure

import (
	"fmt"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/introspect"
	"os"
)

const intro = `
<node>
	<interface name="org.mpris.MediaPlayer2.Player">
		<method name="PlayPause"/>
		<method name="Next"/>
		<method name="Previous"/>
		<method name="Pause"/>
		<method name="Stop"/>
	</interface>` + introspect.IntrospectDataString + `</node> `

type DbusMethods struct {
	publisher *sconsify.Publisher
}

func (dbus DbusMethods) PlayPause() *dbus.Error {
	dbus.publisher.PlayPauseToggle()
	return nil
}

func (dbus DbusMethods) Next() *dbus.Error {
	dbus.publisher.NextPlay()
	return nil
}

func (dbus DbusMethods) Previous() *dbus.Error {
	//dbus.publisher.Previous()
	return nil
}

func (dbus DbusMethods) Pause() *dbus.Error {
	dbus.publisher.Pause()
	return nil
}

func (dbus DbusMethods) Stop() *dbus.Error {
	//dbus.publisher.Stop()
	return nil
}

func StartDbus(publisher *sconsify.Publisher) {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot access dbus, ignoring...")
		return
	}
	reply, err := conn.RequestName("org.mpris.MediaPlayer2.sconsify", dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot request dbus name, ignoring...")
		return
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "org.mpris.MediaPlayer2.sconsify name already taken, ignoring...")
		return
	}
	dbusMethods := new(DbusMethods)
	dbusMethods.publisher = publisher
	conn.Export(dbusMethods, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")
	conn.Export(introspect.Introspectable(intro), "/org/mpris/MediaPlayer2", "org.freedesktop.DBus.Introspectable")
	//select {}
}
