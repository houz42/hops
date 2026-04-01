Output demo/add-profile.webp

Set Shell "bash"
Set FontFamily "JetBrainsMono Nerd Font Mono"
Set FontSize 18
Set LetterSpacing 0
Set LineHeight 1.2
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 80ms
Set Framerate 30
Set WindowBar Colorful
Set Margin 20
Set MarginFill "#1e1e2e"
Set BorderRadius 10

# Launch hops
Type "./hops -f /tmp/hops-demo/hosts"
Enter
Sleep 2s

# Start adding a profile
Type "a"
Sleep 1s

# Type profile name
Type "dev-local"
Enter
Sleep 1s

# Add entries
Type "127.0.0.1 myapp.local"
Enter
Sleep 800ms

Type "127.0.0.1 api.myapp.local"
Enter
Sleep 800ms

Type "127.0.0.1 admin.myapp.local"
Enter
Sleep 800ms

# Empty enter to finish
Enter
Sleep 1.5s

# See it in the list, enable it
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Enter
Sleep 1.5s

# Quit
Type "q"
Sleep 1s
