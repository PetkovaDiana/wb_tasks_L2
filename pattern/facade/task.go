package

import "fmt"

/*
Паттерн «фасад»
*/

type TV struct{}

func (t *TV) On() {
	fmt.Println("TV is On")
}

func (t *TV) Off() {
	fmt.Println("TV is Off")
}

type SoundSystem struct{}

func (s *SoundSystem) On() {
	fmt.Println("Sound System is On")
}

func (s *SoundSystem) Off() {
	fmt.Println("Sound System is Off")
}

func (s *SoundSystem) SetVolumeLevel(level int) {
	fmt.Printf("Sound system volume set to %d\n", level)
}

type DVDPlayer struct{}

func (d *DVDPlayer) On() {
	fmt.Println("DVD Player is ON")
}

func (d *DVDPlayer) Off() {
	fmt.Println("DVD Player is OFF")
}

func (d *DVDPlayer) Play(movie string) {
	fmt.Printf("Playing movie: %s\n", movie)
}

// HomeTheaterFacade - структура для управления
type HomeTheaterFacade struct {
	tv          *TV
	soundSystem *SoundSystem
	dvdPlayer   *DVDPlayer
}

func NewHomeTheaterFacade(tv *TV, soundSystem *SoundSystem, dvdPlayer *DVDPlayer) *HomeTheaterFacade {
	return &HomeTheaterFacade{
		tv:          tv,
		soundSystem: soundSystem,
		dvdPlayer:   dvdPlayer,
	}
}

func (h *HomeTheaterFacade) WatchMovie(movie string) {
	fmt.Println("Get ready to watch a movie...")
	h.tv.On()
	h.soundSystem.On()
	h.soundSystem.SetVolumeLevel(5)
	h.dvdPlayer.On()
	h.dvdPlayer.Play(movie)
}

func (h *HomeTheaterFacade) EndMovie() {
	fmt.Println("The movie is over...")
	h.tv.Off()
	h.soundSystem.Off()
	h.dvdPlayer.Off()
}
func main() {
	tv := &TV{}
	soundSystem := &SoundSystem{}
	dvdPlayer := &DVDPlayer{}

	homeTheater := NewHomeTheaterFacade(tv, soundSystem, dvdPlayer)

	homeTheater.WatchMovie("Harry Potter")
	homeTheater.EndMovie()
}
