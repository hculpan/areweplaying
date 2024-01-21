# Are We Playing?
A Go web app to track if players are going to attend the next session of my TTPRG.

The purpose of this app is for my players to indicate if they'll attend the game
with as little effort as possible. To this end, I didn't want to have a site that
they'd have to log into, and I'd prefer they don't have to even go to the website
unless they want to. So players can interace with this app using email only.

It consists of two components, a web app and a CLI. For players, the interface is
very simple. It just presents a list of players, with an indication of whether or
not they will attend. They can click to update this to say that they plan to attend
or that they will not attend. There is no security here, as the only thing that 
anyone can do is update attendance. No personal information, other than player names
(presumably first names) will be displayed or updated here. Players are configured
manually by me by directly editing the json file that is the only data store for
this project.

This web app does have a "Setup" screen that can do some admin, like send a
game reminder email to all the players, or advance to the next session date. This
does require a password to login (though not a username), and I should be the only
one to ever see this. The security on this app is by no means ironclad, but it's
pretty standard. Once I login, it saves a JWT token as a cookie and every page 
requires this token.

The second component is a CLI, and it just checks an email inbox to see if any of
the players have responded to the game reminder email with a "yes" or "no", indicating
if they will attend. The CLI runs as a cron job, and if it sees a new email from
one of the player's emails, it looks for a "yes" or "no" in the first non-blank line.
If it finds this, it updates the player's status accordingly.

# Development

This app has a number of dependencies. It uses the Chi router and HTMX for the
web interface with Bootstrap 5 for styling. And it uses ```templ``` for templating
the HTML.

The data store is just a json file, and an example json file is included. By default,
it stores data to ```players.json``` in the working directory. 

It does use Godotenv for configuration, and Cobra for the CLI.

The only other significant library is ```go-imap``` to simplify interfacing with
IMAP and SMTP servers.

I also used Air during development for hot restarts, so there's an ```.air.toml```
file in the project. This is just for the web app; it is not configured to work with
the CLI app.

# Building

```make build```

This will produce two executables:
```areweplaying``` - The web app
```awp-cli``` - The CLI app, obviously

I've also included a command to make Linux/AMD64 versions of the files:

```make linuxbuild```

This will produce the same two executables, but with a ```.linux``` extension.