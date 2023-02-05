package web

import (
	"embed"
)

//go:embed *
//go:embed _app/version.json
//go:embed _app/immutable/assets/_layout-d5b06b35.css
//go:embed _app/immutable/chunks/*
//go:embed _app/immutable/components/error.svelte-9a87daf8.js
//go:embed _app/immutable/components/pages/_layout.svelte-a9e68a06.js
//go:embed _app/immutable/components/pages/_page.svelte-099640e0.js
//go:embed _app/immutable/components/pages/battles/*
//go:embed _app/immutable/components/pages/leaderboards/*
//go:embed _app/immutable/components/pages/monsters/*
//go:embed _app/immutable/components/pages/playground/*
//go:embed _app/immutable/modules/pages/_layout.ts-9cbb603b.js
var Resource embed.FS
