# nice shot

This GUI is about visulizing the data your coffee shots are producing when using you 
belove coffee machine. It helps you investigate if your most beloved drink became 
perfect without tasting, but by analysing it.

Enjoy!

## Nice Shot – Full‑stack demo



This project contains:

- Backend: Go (Echo) serving mocked espresso shot data and statistics.
- Frontend: Vue 3 (Vite) + TailwindCSS visualizing trends and recent shots.

### Requirements

- Go 1.23+ (toolchain will auto-select go1.24 on first build)
- Node 20+ and npm

### Run backend (Echo)

```bash
cd backend
go run .
```

Server runs at `http://localhost:8080` with these routes:

- `GET /api/health`
- `GET /api/shots?limit=100` – most recent first

### Run frontend (Vue + Vite)

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173`. The frontend expects the backend at `http://localhost:8080/api`. To change it, set `VITE_API_BASE` in an `.env` file in `frontend/`.

### Notes

- Mock dataset is generated at backend startup with ~250 shots between Aug 1 and today.
- UI includes: daily brew time line chart, grind vs peak pressure scatter, and a sortable table of recent shots.

