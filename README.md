# unitkit

> Agentic-first unit conversion service. Convert between units of length, mass, temperature, volume, area, speed, time, data storage, pressure, energy, power, angle, and frequency. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
make build
./unitkit
# unitkit listening on :8080
```

## API

### Convert (GET)

```bash
curl "http://localhost:8080/convert?from=m&to=ft&value=1"
# from=m to=ft value=1 result=3.280839895

curl "http://localhost:8080/convert?from=C&to=F&value=25"
# from=C to=F value=25 result=77
```

### Convert (POST)

```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"from":"kg","to":"lb","value":1}'
# from=kg to=lb value=1 result=2.204622622
```

### Convert (Path)

```bash
curl "http://localhost:8080/convert/m/km/1000"
# from=m to=km value=1000 result=1
```

### List Categories

```bash
curl "http://localhost:8080/categories"
# name=angle description=Angle measurements base=deg units=deg,rad,grad,...
# name=area description=Area measurements base=m2 units=m2,km2,...
# ...
```

### List Units in a Category

```bash
curl "http://localhost:8080/units/length"
# category=length base=m
# symbol=m name=meter
# symbol=km name=kilometer
# ...
```

### JSON Responses

Add `Accept: application/json` header or `?format=json` query param:

```bash
curl "http://localhost:8080/convert?from=m&to=ft&value=1&format=json"
# {"from":"m","to":"ft","value":1,"result":3.280839895,"category":"length"}
```

### Self-Documenting

```bash
curl "http://localhost:8080/help"
# Full operating manual for AI agents
```

### MCP Endpoint

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","id":1}'
```

## Categories

| Category | Units |
|----------|-------|
| length | m, km, cm, mm, um, nm, mi, yd, ft, in, nmi, ly, au, pc |
| mass | kg, g, mg, ug, t, lb, oz, st, ton_us, ton_uk, ct, gr |
| temperature | C, F, K, R |
| volume | L, mL, m3, cm3, mm3, ft3, in3, gal_us, gal_uk, qt_us, pt_us, cup_us, floz_us, tbsp_us, tsp_us, bbl_oil |
| area | m2, km2, cm2, mm2, ha, ac, ft2, in2, yd2, mi2, nmi2 |
| speed | m/s, km/h, mph, ft/s, kn, mach, c |
| time | ns, us, ms, s, min, h, d, wk, mo, yr, dec, cent |
| data | bit, B, KB, MB, GB, TB, PB, KiB, MiB, GiB, TiB, PiB, Kb, Mb, Gb |
| pressure | Pa, kPa, MPa, GPa, bar, mbar, atm, psi, mmHg, inHg, torr, hPa |
| energy | J, kJ, MJ, cal, kcal, Wh, kWh, MWh, BTU, ftlb, eV, erg, therm |
| power | W, kW, MW, GW, mW, hp, hp_m, BTU/h, ftlb/s, cal/s |
| angle | deg, rad, grad, arcmin, arcsec, turn |
| frequency | Hz, kHz, MHz, GHz, THz, rpm, rad/s |

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `UNITKIT_ADDR` | `:8080` | Listen address |

## Build

```bash
make build    # CGO_ENABLED=0, static binary
make test     # go test -race
make vet      # go vet
```

## License

MIT
