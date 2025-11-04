# Emacs remote integration

This project provides remote Emacs integration.  Basically, if you
want to run Emacs on a remote server via SSH, you can integrate it
with your local machine's clipboard and web browser.

There are two parts:

* `reintegrate.el`, an Emacs package
* `emacs-integration`, an HTTP server

`reintegrate.el` in the remote Emacs sends HTTP requests to
`emacs-integration` running on your local machine through a forwarded
port.

## Quickstart

On your local machine:

1. Install clipboard utilities:

   * For X, install xsel.
   * For Wayland, install wl-clipboard.

2. Build the integration server (requires Go):

        (cd emacs-integration && go install .)  # Installs into $HOME/go/bin

3. Assuming your local machine uses systemd, copy the service and
   socket files into `~/.config/systemd/user`:

         cp *.socket *.service ~/.config/systemd/user

4. If you're using Wayland, edit the service file and replace the `ExecStart`
   line with the commented one.

5. Enable systemd socket activation:

         systemctl --user enable --now emacs-integration.socket

You don't have to start the service manually, as systemd can do it
automatically when something conects to the port.

On your remote machine:

1. Add `reintegrate.el` to your `load-path`.
2. Add to your `init.el`:

        (when (getenv "SSH_TTY")
          (require 'reintegrate)
          (reintegrate))

When using:

1. SSH into your remote machine and forward the ports:

        ssh -R 9999:localhost:9999 myhost

2. Start Emacs.
