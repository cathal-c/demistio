# catalog

The `catalog` package provides a fluent, programmatic API for constructing Istio configuration resources (`config.Config`) in a safe, readable, and testable way.

It acts as a builder-layer abstraction on top of Istio’s generic `config.Config` structure, which wraps protobuf-based resource specifications such as `DestinationRule`, `ServiceEntry`, and `PeerAuthentication`.

## 📦 Purpose

Istio’s internal config system is intentionally generic (`Spec any`), making it easy to use but tedious to construct correctly. The `catalog` package solves this by offering helper builders that:

- Ensure the correct `GroupVersionKind` is set
- Type-assert and initialize the correct `Spec` type
- Provide a fluent interface for chaining and mutation
- Keep your test and setup code concise and expressive

## ✅ Supported Resources

- [`DestinationRule`](https://istio.io/latest/docs/reference/config/networking/destination-rule/)
- [`ServiceEntry`](https://istio.io/latest/docs/reference/config/networking/service-entry/)
- [`PeerAuthentication`](https://istio.io/latest/docs/reference/config/security/peer_authentication/)

## 🔨 Example Usage

```go
import "github.com/cathal-c/demistio/pkg/catalog"
import v1 "istio.io/api/networking/v1"

// Construct a strict mTLS PeerAuthentication
pa := catalog.NewPeerAuthenticationBuilder("default", "istio-system").
    WithMutualTlsMode(v1.PeerAuthentication_MutualTLS_STRICT).
    Build()

// Construct a DestinationRule with ISTIO_MUTUAL TLS
dr := catalog.NewDestinationRuleBuilder("default", "istio-system").
    WithHost("*.svc.cluster.local").
    WithTlsMode(v1.ClientTLSSettings_ISTIO_MUTUAL).
    Build()

// Construct a basic ServiceEntry
se := catalog.NewServiceEntryBuilder("google", "istio-system").
    WithHosts("google.com").
    Build()

// You can now add these directly to Istio's config store
