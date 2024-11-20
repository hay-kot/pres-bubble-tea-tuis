---
# You can also start simply with 'default'
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: https://cover.sli.dev
# some information about your slides (markdown enabled)
title: Building TUIs in Go
info: |
  ## Bubble Tea TUIs

  How to build Terminal UIs Using Bubble Tea

class: text-center
drawings:
  persist: false
transition: slide-left
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
overviewSnapshots: true
dragPos:
  haykot-logo: 17,476,61,62
---

# Building Terminal User Interfaces in Go

How to build interactive terminal UIs in Go using Bubble Tea

<div class="abs-br m-6 flex gap-2">
  <button @click="$slidev.nav.openInEditor()" title="Open in Editor" class="text-xl slidev-icon-btn opacity-50 !border-none !hover:text-white">
    <carbon:edit />
  </button>
  <a href="https://github.com/hay-kot/pres-bubble-tea-tuis" target="_blank" alt="GitHub" title="Open in GitHub"
    class="text-xl slidev-icon-btn opacity-50 !border-none !hover:text-white">
    <carbon-logo-github />
  </a>
</div>

<a v-drag="'haykot-logo'" href="https://haykot.dev" target="_blank">
  <img src="/images/haykot-logo.jpg"
</a>

---
dragPos:
    headshot: 765,39,151,149
---

# About Me

- Senior Software Engineer (Johnson Controls)
- Mostly writing C# day-to-day
- Written applications in Go, C#, Python, JavaScript, and Rust
- Launched a Sass Application
  - [Recipinned](https://recipinned.com) - Recipinned allows you to import Recipes from the web, access step-by-step cooking instructions, plan and shop for your meals, and create your one of a kind cookbook.
- Open Source Contributor [github.com/hay-kot](https://github.com/hay-kot)
  - Mealie (Recipe Manager) 10m+ downloads
  - Homebox (Inventory Manager) 1m+ downloads
- Originally From Alaska, moved to Minnesota in summer of 2023
- I love cycling, running, and skiing

<img src="/images/headshot.jpg" class="rounded-full" v-drag="'headshot'" />

---
---

# What's a TUI???

- Interactive Terminal Application
- Runs in your Terminal
- Contains some kind of render loop to refresh the page
- Allows for interactions through keyboard and/or mouse.
- Example Applications
  - [Lazygit](https://github.com/jesseduffield/lazygit) - Git Workflows
  - [Lazydocker](https://github.com/jesseduffield/lazydocker) - Docker Management
  - [Posting](https://github.com/darrenburns/posting) - Postman Alternative

---
layout: image
image: "/images/lazygit-demo.gif"
---

---
dragPos:
    bubbletea: 615,325,240,153
---

# Bubble Tea and Go

- Bubble Tea is a Framework for Build TUIs in Go
- Ships with several companion libraries
  - [Bubbles](https://github.com/charmbracelet/bubbles) - Common Component Library
  - [Lipgloss](https://github.com/charmbracelet/lipgloss) - Css Like Style Framekwork
  - [Harmonica](https://github.com/charmbracelet/harmonica) - A simple, physics-based animation library
- Checkout [charm.sh](https://charm.sh/) for all their projects

<img src="/images/bubble-tea.png" v-drag="'bubbletea'" />

---
layout: section
---

# Covering The Basics

---
---

# The Basics: Understanding The Model

```go
type Model interface {
	// Init is the first function that will be called. It returns an optional
	// initial command. To not perform an initial command return nil.
	Init() Cmd

	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the model and/or send a command.
	Update(Msg) (Model, Cmd)

	// View renders the program's UI, which is just a string. The view is
	// rendered after every Update.
	View() string
}
```

<br/>

This is the core type that you'll use to create your interactive elements

---
---

# The Basics: The Render Loop

<img src="/images/renderloop.png"/>

---
---

# The Basics: Commands and Messages

```go
// Cmd is an IO operation that returns a message when it's complete. If it's
// nil it's considered a no-op. Use it for things like HTTP requests, timers,
// saving and loading from disk, and so on.
//
// Note that there's almost never a reason to use a command to send a message
// to another part of your program. That can almost always be done in the
// update function.
type Cmd func() Msg

// Msg contain data from the result of a IO operation. Msgs trigger the update
// function and, henceforth, the UI.
type Msg interface{}
```

<br/>

`Cmd` and `Msg` are the core types to create and dispatch events within the _"Loop"_


---
---



# The Basics: Running The Model

````md magic-move {lines: true}
```go
func TUI(ctx context.Context) error {
	p := tea.NewProgram(pages.NewLayout(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}
```

```go {2}
func TUI(ctx context.Context) error {
	p := tea.NewProgram(pages.NewLayout(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}
```

````

<br/>

- `tea.NewProgram` takes in a `Model` interface and starts the rendering loop.
- `tea.WithAltScreen()` runs the application in 'take over mode' - it will take up the whole terminal screen.

---
---

# Live Demo

- Overview of Existing Project Structure
- Debugging advice
- Basic model rendering
- Creating "Page" like views
- Running async tasks

---
layout: two-cols
dragPos:
  presentation: 492,280,179,166
  qrcode: 57,280,179,166
---

### Linkedin

[www.linkedin.com/in/hay-kot/](https://linkedin.com/in/hay-kot/)

<br/>
<br/>

### My Website

[haykot.dev](https://haykot.dev)

<img src="/images/qr-haykot.dev.png" v-drag="`qrcode`"/>


::right::

### Bluesky

[bsky.app/profile/haykot.dev](https://bsky.app/profile/haykot.dev)

<br/>
<br/>

### This Presentation

[hay-kot.github.io/pres-bubble-tea-tuis](https://hay-kot.github.io/pres-bubble-tea-tuis/1)

<img src="/images/qr-presentation.png" v-drag="`presentation`"/>
