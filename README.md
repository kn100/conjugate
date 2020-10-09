# Conjugate - Convert Youtube Music Links into Spotify links

A CLI thing that when given a Youtube Music link, it gives you the closest match from Spotify.

This exists because I was tired of being one of the few Youtube Music users in the office where I work.

## Getting it

Go check https://github.com/kn100/conjugate/releases and download the release that best suits your system!

## Usage

After you've copied the `conjugate` binary somewhere in your path, run `conjugate -reconfigure`. It will request some configuration details. You'll need to provide it a Youtube Data API key, and a Spotify Client ID and secret.

From then on, you can run `conjugate` with the `-y=<link to song on Youtube or Youtube music>` flag, and if a match is found on Spotify, the link will be returned.

Also provided is a `-raw` flag - which when provided will output the found link with no formatting. This is so you can pipe the output elsewhere easily.

Try it now with this example command `conjugate -y="https://www.youtube.com/watch?v=dQw4w9WgXcQ"`.
