package fake

import "fmt"

// Installer only for test purpose
type Installer struct {
	hasError bool
	support  bool
}

// NewFakeInstaller returns a Installer
func NewFakeInstaller(support bool, hasError bool) *Installer {
	return &Installer{
		hasError: hasError,
		support:  support,
	}
}

// Available check if support current platform
func (d *Installer) Available() bool {
	return true
}

// Install installs the vim
func (d *Installer) Install() (err error) {
	if d.hasError {
		err = fmt.Errorf("fake error")
	}
	return
}

// Uninstall uninstalls the vim
func (d *Installer) Uninstall() (err error) {
	if d.hasError {
		err = fmt.Errorf("fake error")
	}
	return
}

// WaitForStart waits for the service be started
func (d *Installer) WaitForStart() (ok bool, err error) {
	if d.hasError {
		err = fmt.Errorf("fake error")
	}
	ok = true
	return
}

// Start starts the vim service
func (d *Installer) Start() (err error) {
	if d.hasError {
		err = fmt.Errorf("fake error")
	}
	return
}

// Stop stops the vim service
func (d *Installer) Stop() (err error) {
	if d.hasError {
		err = fmt.Errorf("fake error")
	}
	return
}
