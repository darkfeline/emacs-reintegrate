;;; reintegrate.el --- Remote integration                    -*- lexical-binding: t; -*-

;; Copyright (C) 2018  Allen Li

;; Author: Allen Li <darkfeline@felesatra.moe>
;; Keywords: comm

;; This program is free software; you can redistribute it and/or modify
;; it under the terms of the GNU General Public License as published by
;; the Free Software Foundation, either version 3 of the License, or
;; (at your option) any later version.

;; This program is distributed in the hope that it will be useful,
;; but WITHOUT ANY WARRANTY; without even the implied warranty of
;; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
;; GNU General Public License for more details.

;; You should have received a copy of the GNU General Public License
;; along with this program.  If not, see <https://www.gnu.org/licenses/>.

;;; Commentary:

;; The package provides integration when running Emacs remotely, for
;; example integrating with your local clipboard and web browser.

;;; Code:

(require 'url)

(defgroup reintegrate nil
  "reintegrate customization group."
  :group 'communication)

(defconst reintegrate-encoding 'utf-8
  "Encoding to use for communication.")

(defcustom reintegrate-host "127.0.0.1:9999"
  "Integration server."
  :type 'string)

;;;###autoload
(defun reintegrate ()
  "Enable remote integration."
  (interactive)
  (setq browse-url-browser-function #'reintegrate-browse-url
        interprogram-cut-function #'reintegrate-cut
        interprogram-paste-function #'reintegrate-paste))

(defun reintegrate-browse-url (url &optional _new-window)
  "Browse URL through Emacs proxy.
NEW-WINDOW is ignored."
  (let ((url-request-method "POST")
        (url-request-extra-headers '(("Content-Type" . "text/plain; charset=utf-8")))
        (url-request-data (encode-coding-string url reintegrate-encoding t)))
    (url-retrieve-synchronously (reintegrate--rpc-path reintegrate-host "/browser") t t)))

(defun reintegrate-cut (string)
  "Make STRING available for pasting through integration server."
  (let ((url-request-method "PUT")
        (url-request-extra-headers '(("Content-Type" . "text/plain; charset=utf-8")))
        (url-request-data (encode-coding-string string reintegrate-encoding t)))
    (url-retrieve-synchronously (reintegrate--rpc-path reintegrate-host "/clipboard") t t)))

(defun reintegrate-paste ()
  "Get string for pasting through integration server."
  (with-current-buffer (url-retrieve-synchronously (reintegrate--rpc-path reintegrate-host "/clipboard") t t)
    (goto-char (point-min))
    (let ((eol (save-excursion (forward-line) (1- (point)))))
      (if (not (search-forward "200 OK" eol t))
          (message "Error getting clipboard: %s" (buffer-substring-no-properties (point) eol))
        (search-forward "\n\n")
        (decode-coding-string (buffer-substring-no-properties (point) (point-max)) reintegrate-encoding)))))

(defun reintegrate--rpc-path (host path)
  "Get URL for RPC at HOST and PATH."
  (format "http://%s%s" host path))

(provide 'reintegrate)
;;; reintegrate.el ends here
