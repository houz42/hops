Output demo/demo.webp

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

# Navigate profiles
Sleep 500ms
Down
Sleep 800ms
Down
Sleep 800ms
Up
Sleep 800ms

# Toggle staging on
Down
Sleep 500ms
Enter
Sleep 1s

# Toggle it back off
Enter
Sleep 1s

# View work profile entries
Up
Sleep 500ms
Right
Sleep 2s

# Navigate entries
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms

# Go back
Left
Sleep 1.5s

# Quit
Type "q"
Sleep 1s
