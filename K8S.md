# Kubernetes Production Readiness

## Scaling
kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog

All 3 pods are running (see screenshot).

## Health Checks

### Readiness vs Liveness Probe

**Readiness Probe:**
Prueft ob der Container bereit ist Traffic zu empfangen.
Wenn sie fehlschlaegt wird der Pod aus dem Service-Load-Balancer entfernt — er bekommt keinen Traffic mehr, wird aber nicht neu gestartet.

**Liveness Probe:**
Prueft ob der Container noch lebt/funktioniert.
Wenn sie fehlschlaegt wird der Pod neu gestartet.

**Unterschiedliche initialDelaySeconds:**
Readiness (5s): Der Container muss schnell bereit sein um Traffic zu empfangen.
Liveness (15s): Der Container braucht mehr Zeit zum Starten bevor Kubernetes entscheidet ob er neu gestartet werden soll.

## Resource Limits

**Was passiert wenn Memory-Limit ueberschritten wird:**
Der Container wird mit OOMKilled (Out of Memory) beendet und neu gestartet.

**Was passiert wenn CPU-Limit ueberschritten wird:**
Der Container wird gedrosselt (throttled) — er laeuft langsamer aber wird nicht neu gestartet.

**Warum requests und limits angeben:**
- requests: Garantiert dem Container diese Ressourcen — Kubernetes verwendet sie fuer das Scheduling
- limits: Verhindert dass ein Container alle Ressourcen des Nodes verbraucht und andere Container beeintraechtigt