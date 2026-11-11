// Code generated from the Status302 public OpenAPI spec (https://status302.com/docs/api/openapi.yaml). DO NOT EDIT.

package status302

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// CreateStreamTicket calls POST /v1/stream/ticket: Ticket for a browser WebSocket connection.
//
// Returns a single-use ticket, valid for 60 s, for opening the realtime
// stream from a browser: `GET /v1/stream?ticket=<ticket>` (WebSocket,
// protocol: https://status302.com/docs/stream/). Long-lived credentials never appear
// in URLs. The connection acts as the caller; ending the session (or
// revoking the API key) closes it with 4401. With an API key (backends
// handing tickets to their own browser clients) the key needs the
// `stream` scope; server-side clients can instead send the key in
// `Authorization: Bearer` on the upgrade itself.
//
// **API key scope:** `stream`.
func (c *Client) CreateStreamTicket(ctx context.Context, opts ...RequestOption) (*StreamTicket, error) {
	var out StreamTicket
	if err := c.do(ctx, http.MethodPost, "/v1/stream/ticket", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventsParams holds the query parameters of ListEvents.
type ListEventsParams struct {
	// Name search (substring or trigram word similarity on the normalized event name). An event ID or event page URL also matches the events it names: a Ticketmaster event ID (16 hex digits, e.g. `380064AEC95E177C`, or a ticketmaster.* `…/event/{id}` URL) or a StubHub/viagogo event number of a linked marketplace listing (e.g. `161219435`, `E-161219435`, or a stubhub.* / viagogo.* `…/event/{n}` or `…/E-{n}` URL).
	Q *string
	// Only events starting at or after this instant. Defaults to now, in which case undated (TBA/TBD) events are included too; pass an earlier instant to see events that already started.
	From *time.Time
	// Only events starting before this instant.
	To *time.Time
	// Venue country (ISO 3166-1 alpha-2).
	Country *string
	// Region code (`GET /v1/regions`); only venues in one of its countries. With `country`, both apply.
	Region *string
	// Venue city (case-insensitive exact match).
	City *string
	// `lat,lng,radiusKm`: only events at venues within `radiusKm` (at most 500) of
	// the point, e.g. `51.5072,-0.1276,25`.
	Near *string
	// Comma-separated classification segments (any of them, case-insensitive), e.g. `Music,Sports`. See `GET /v1/catalog/facets`.
	Segment []string
	// Comma-separated genres (any of them, case-insensitive).
	Genre []string
	// Comma-separated lifecycle states to include (see `EventLifecycle`).
	Lifecycle []EventLifecycle
	// Comma-separated polling tiers to include (see `EventTier`).
	Tier []EventTier
	// Only events with (`true`) or without (`false`) a compiled seat map.
	HasMap *bool
	// Comma-separated sources the event must be available at, all of them:
	// `tm` (hosted on Ticketmaster), `dice` (a Dice event) and marketplaces
	// (`stubhub` — StubHub & viagogo, `crowdvolt`, `tickpick`, `vividseats`, `ticombo` — Ticombo, `tmresale` — Ticketmaster resale listings, the `tm_resale` market) with a linked listing.
	Sources []string
	// Comma-separated statuses to include.
	Status []EventStatus
	// Only events of this primary seller (`EventSummary.primarySource`): `tm` for events hosted on Ticketmaster, `dice` for Dice's own events.
	PrimarySource *PrimarySource
	// Only events hosted on Ticketmaster (with a `tmEventId`). Same as `primarySource=tm`, which replaces it; 400 together with another `primarySource`.
	TmOnly *bool
	// Only events first seen by the catalog at or after this instant.
	NewSince *time.Time
	// Only events at this venue.
	VenueID *int64
	// Only events featuring this attraction.
	AttractionID *int64
	// Market scope (default `primary`).
	Scope *MetricScope
	// Minimum floor of the scope (all-in minor units); excludes events without a floor.
	MinFloor *int64
	MaxFloor *int64
	// Minimum 24h floor change (fraction, 0.1 = +10%); excludes events without a baseline.
	MinChg24h *float64
	MaxChg24h *float64
	// Minimum tickets sold in the last 24h.
	MinSales24h *int
	// Minimum spread of the scope (fraction, 0.2 = the floor is 20% above the primary floor); marketplace scopes carry a spread. Excludes events without one.
	MinSpread *float64
	// Maximum primary sell-out estimate in hours (`EventMetrics.selloutHours`); only the `primary` scope carries one. Excludes events without one.
	MaxSelloutHours *float64
	// Only events whose public on-sale started at or after this instant (excludes events without one).
	OnsaleFrom *time.Time
	// Only events whose public on-sale starts before this instant (excludes events without one).
	OnsaleTo *time.Time
	// Comma-separated event IDs (at most 200, more is a 400), e.g. `123,456`. Only these events; combines with every other filter, the sort and `scope`. Unknown IDs match nothing.
	IDs []int64
	// Only the events in this saved list of the active organization (404 when there is none); combines with `ids` and every other filter.
	ListID *string
	// `startsAt` (ascending, undated events last), `firstSeen` (newest first),
	// `onsale` (public on-sale start, latest first; events without one excluded), or
	// any numeric metric of `scope` (the `EventMetrics` field of the same name;
	// events without the value last).
	Sort *ListEventsSort
	// Direction of a metric sort (default `asc` for `floor` and `selloutHours`, `desc` otherwise); not allowed with `startsAt`/`firstSeen`/`onsale`.
	Order *ListEventsOrder
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListEventsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Q != nil {
		q.Set("q", *p.Q)
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	if p.Country != nil {
		q.Set("country", *p.Country)
	}
	if p.Region != nil {
		q.Set("region", *p.Region)
	}
	if p.City != nil {
		q.Set("city", *p.City)
	}
	if p.Near != nil {
		q.Set("near", *p.Near)
	}
	if len(p.Segment) > 0 {
		q.Set("segment", joinQuery(p.Segment, func(v string) string { return v }))
	}
	if len(p.Genre) > 0 {
		q.Set("genre", joinQuery(p.Genre, func(v string) string { return v }))
	}
	if len(p.Lifecycle) > 0 {
		q.Set("lifecycle", joinQuery(p.Lifecycle, func(v EventLifecycle) string { return string(v) }))
	}
	if len(p.Tier) > 0 {
		q.Set("tier", joinQuery(p.Tier, func(v EventTier) string { return string(v) }))
	}
	if p.HasMap != nil {
		q.Set("hasMap", strconv.FormatBool(*p.HasMap))
	}
	if len(p.Sources) > 0 {
		q.Set("sources", joinQuery(p.Sources, func(v string) string { return v }))
	}
	if len(p.Status) > 0 {
		q.Set("status", joinQuery(p.Status, func(v EventStatus) string { return string(v) }))
	}
	if p.PrimarySource != nil {
		q.Set("primarySource", string(*p.PrimarySource))
	}
	if p.TmOnly != nil {
		q.Set("tmOnly", strconv.FormatBool(*p.TmOnly))
	}
	if p.NewSince != nil {
		q.Set("newSince", (*p.NewSince).Format(time.RFC3339Nano))
	}
	if p.VenueID != nil {
		q.Set("venueId", strconv.FormatInt(*p.VenueID, 10))
	}
	if p.AttractionID != nil {
		q.Set("attractionId", strconv.FormatInt(*p.AttractionID, 10))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.MinFloor != nil {
		q.Set("minFloor", strconv.FormatInt(*p.MinFloor, 10))
	}
	if p.MaxFloor != nil {
		q.Set("maxFloor", strconv.FormatInt(*p.MaxFloor, 10))
	}
	if p.MinChg24h != nil {
		q.Set("minChg24h", strconv.FormatFloat(*p.MinChg24h, 'f', -1, 64))
	}
	if p.MaxChg24h != nil {
		q.Set("maxChg24h", strconv.FormatFloat(*p.MaxChg24h, 'f', -1, 64))
	}
	if p.MinSales24h != nil {
		q.Set("minSales24h", strconv.Itoa(*p.MinSales24h))
	}
	if p.MinSpread != nil {
		q.Set("minSpread", strconv.FormatFloat(*p.MinSpread, 'f', -1, 64))
	}
	if p.MaxSelloutHours != nil {
		q.Set("maxSelloutHours", strconv.FormatFloat(*p.MaxSelloutHours, 'f', -1, 64))
	}
	if p.OnsaleFrom != nil {
		q.Set("onsaleFrom", (*p.OnsaleFrom).Format(time.RFC3339Nano))
	}
	if p.OnsaleTo != nil {
		q.Set("onsaleTo", (*p.OnsaleTo).Format(time.RFC3339Nano))
	}
	if len(p.IDs) > 0 {
		q.Set("ids", joinQuery(p.IDs, func(v int64) string { return strconv.FormatInt(v, 10) }))
	}
	if p.ListID != nil {
		q.Set("listId", *p.ListID)
	}
	if p.Sort != nil {
		q.Set("sort", string(*p.Sort))
	}
	if p.Order != nil {
		q.Set("order", string(*p.Order))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListEvents calls GET /v1/events: Search the event catalog.
//
// Lists catalog events with cursor pagination. Pass `nextCursor` from a
// response as `cursor` (with the same filters and sort) to get the next page.
// Each row carries the `metrics` of `scope` (default `primary`) when the
// event has market data. Events that started more than 24 hours ago carry
// no metrics without ended-event history (metric filters and sorts treat
// them as having no value).
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEvents(ctx context.Context, params *ListEventsParams, opts ...RequestOption) (*EventPage, error) {
	var out EventPage
	if err := c.do(ctx, http.MethodGet, "/v1/events", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventsPages iterates the pages of ListEvents, following nextCursor as cursor until the last page.
func (c *Client) ListEventsPages(ctx context.Context, params *ListEventsParams, opts ...RequestOption) iter.Seq2[*EventPage, error] {
	var p ListEventsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*EventPage, error) { return c.ListEvents(ctx, &p, opts...) }, func(page *EventPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListEventsAll iterates every item of ListEvents across pages.
func (c *Client) ListEventsAll(ctx context.Context, params *ListEventsParams, opts ...RequestOption) iter.Seq2[EventSummary, error] {
	return items(c.ListEventsPages(ctx, params, opts...), func(page *EventPage) []EventSummary { return page.Items })
}

// ListRegions calls GET /v1/regions: Regions.
//
// The regions of the `region` filter (catalog, calendar, sales feed and
// summary, `sales`/`onsales` stream topics): named groups of venue
// countries with a display currency, in display order. Regions do not
// overlap.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListRegions(ctx context.Context, opts ...RequestOption) (*RegionList, error) {
	var out RegionList
	if err := c.do(ctx, http.MethodGet, "/v1/regions", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCatalogFacets calls GET /v1/catalog/facets: Segments and genres of the catalog.
//
// The values of the `segment` and `genre` filters of `GET /v1/events` with their numbers of upcoming (or undated) events.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetCatalogFacets(ctx context.Context, opts ...RequestOption) (*CatalogFacets, error) {
	var out CatalogFacets
	if err := c.do(ctx, http.MethodGet, "/v1/catalog/facets", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEvent calls GET /v1/events/{id}: Event detail.
//
// One catalog event including its public on-sale and presale windows.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEvent(ctx context.Context, id int64, opts ...RequestOption) (*EventDetail, error) {
	var out EventDetail
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10)), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventMap calls GET /v1/events/{id}/map: The event's seat map.
//
// Status of the event's seat map and, once compiled, the artifact it uses
// (fetched from the event's own Ticketmaster map, never "latest for the
// venue"). `pending` until the first map fetch; `no_map` when Ticketmaster
// has none (GA-only events); `error` when the map could not be compiled (the
// last good artifact, if any, is still returned).
//
// **API key scope:** `maps:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEventMap(ctx context.Context, id int64, opts ...RequestOption) (*EventMap, error) {
	var out EventMap
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/map", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventMapOverlayParams holds the query parameters of GetEventMapOverlay.
type GetEventMapOverlayParams struct {
	// Add the section layers with this mode's values (none by default).
	Mode *OverlayMode
	// Only section layers of this market (all by default).
	Market *OverlayMarket
	// Only the section layer of this marketplace (implies `market=secondary`).
	Source *string
}

func (p *GetEventMapOverlayParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Mode != nil {
		q.Set("mode", string(*p.Mode))
	}
	if p.Market != nil {
		q.Set("market", string(*p.Market))
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	return q
}

// GetEventMapOverlay calls GET /v1/events/{id}/map/overlay: Live availability and prices on the seat map.
//
// Ticketmaster primary inventory of the event's last accepted snapshot and
// its current TM resale listings (source `tmresale`), keyed by the
// artifact's place and section indexes. Bitsets are base64, one bit per
// place, LSB first (place *i* is bit *i*&nbsp;mod&nbsp;8 of byte
// *i*&nbsp;div&nbsp;8). `priceLevel` is base64 of a little-endian int16 per
// place: an index into `priceLevels`, or -1 when unpriced. Primary seats
// join through the manifest's placeIds, TM resale listings through the
// placeIds they name; GA stock through the section names.
// The ETag changes with the seat map, the snapshot and the pricing; send it
// back in `If-None-Match` to get a 304.
//
// With `mode`, the body also carries `sectionLayers`: per scope (Ticketmaster
// `primary` and `tm_resale`, each marketplace `secondary:<source>` with
// section data and their best `secondary:best`), the values of every section
// on the map, keyed by section index and key: tickets (`availability`), plus
// the median (`price`), 24h/72h sales (`sales`), or the spread and ROI
// against the section's primary floor (`spread`, see
// `GET /v1/events/{id}/sections`). With map prices every mode but
// `availability` also carries each section's floor and its 24h change
// (`delta` carries only those), so a hover can show them in every mode.
// Marketplace layers (`secondary:<source>`, `secondary:best`) also break
// each section's tickets down by the row the listings name (`rows`), with
// the seats of listings that name them, so a map can light the exact
// seats, else the row, else the whole section.
// `mode=availability` also strips the price layers like plans without map
// prices get them. `market` and `source` narrow the layers; marketplace
// layers need a plan with secondary marketplaces (omitted otherwise, 402
// when asked for), and the `price`, `delta` and `spread` modes map prices
// (402). With `mode`, the ETag also follows the section data, and for
// `sales` and every mode carrying the 24h change it changes at least every
// 5 minutes. Section layers of an event that started more than 24 hours
// ago need ended-event history (402). Live updates: stream topic `event.{id}.map`.
//
// **API key scope:** `maps:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEventMapOverlay(ctx context.Context, id int64, params *GetEventMapOverlayParams, opts ...RequestOption) (*MapOverlay, error) {
	var out MapOverlay
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/map/overlay", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventSectionsParams holds the query parameters of ListEventSections.
type ListEventSectionsParams struct {
	// Market scope (default `primary`): `primary`, `tm_resale`, `secondary:best` or
	// `secondary:<marketplace>` (any marketplace id, e.g. `secondary:stubhub`).
	Scope *string
	// Share of the resale price the seller pays, for `roi` (default 0).
	SellerFee *float64
}

func (p *ListEventSectionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", *p.Scope)
	}
	if p.SellerFee != nil {
		q.Set("sellerFee", strconv.FormatFloat(*p.SellerFee, 'f', -1, 64))
	}
	return q
}

// ListEventSections calls GET /v1/events/{id}/sections: Section table of an event.
//
// Every section of the event with stock in the scope now or sales in the last
// 72 hours: current tickets, floor and median (all-in per ticket, ticket-weighted
// median), when they last changed, the 24h floor change, tickets sold in 24h and
// 72h with the 24h average sale price, the spread against the primary floor and
// an estimated ROI. Section keys are the event's seat map section names (`index`
// is the artifact section index); Ticketmaster sections no map section matched
// keep their manifest name. The primary floor of a section is its current one
// or, once sold out, the last one Ticketmaster showed; the spread compares the
// scope's floor with it (`primary`: the best marketplace floor), the ROI the
// median resale price over 72h (`primary`: marketplace sales) after `sellerFee`.
// Prices are in `currency` (other currencies converted at the day's rate).
// Secondary scopes need a plan with secondary marketplaces. Plans without map
// prices (`mapPrices`) get availability only, like the map overlay: tickets
// and sales counts, with `currency`, `floorMinor`, `medianMinor`,
// `saleAvg24hMinor`, `floorChg24h`, `spread` and `roi` null. Events that
// started more than 24 hours ago need ended-event history (402). Live section
// values: stream topic `event.{id}.map`.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventSections(ctx context.Context, id int64, params *ListEventSectionsParams, opts ...RequestOption) (*EventSections, error) {
	var out EventSections
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/sections", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventSectionSeriesParams holds the query parameters of GetEventSectionSeries.
type GetEventSectionSeriesParams struct {
	// Market scope (default `primary`): `primary`, `tm_resale`, `secondary:best` or
	// `secondary:<marketplace>` (any marketplace id, e.g. `secondary:stubhub`).
	Scope *string
	Res   *SeriesResolution
	// Range start (default `to` − 72h).
	From *time.Time
	// Range end, exclusive (default now).
	To *time.Time
}

func (p *GetEventSectionSeriesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", *p.Scope)
	}
	if p.Res != nil {
		q.Set("res", string(*p.Res))
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	return q
}

// GetEventSectionSeries calls GET /v1/events/{id}/sections/{key}/series: Chart series of one section.
//
// Column-oriented buckets of one section's history in a scope: floor, median and
// tickets holding at the end of each bucket (each change holds until the next;
// null before the first known value, floor and median null while nothing is
// priced) and the tickets sold in it. `secondary:best` takes the lowest floor and
// the summed tickets across marketplaces. Without `res` the resolution follows
// the range: 5m up to 3 days, 1h up to 60 days, 1d beyond; more than 2000 buckets
// is a 400. History reaches back as far as the plan's chart history. Plans
// without map prices (`mapPrices`) get availability only: `currency` null and
// every `floor` and `median` bucket null; `tickets` and `sales` stay.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEventSectionSeries(ctx context.Context, id int64, key string, params *GetEventSectionSeriesParams, opts ...RequestOption) (*SectionSeries, error) {
	var out SectionSeries
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/sections/"+url.PathEscape(key)+"/series", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListInventoryChangesParams holds the query parameters of ListInventoryChanges.
type ListInventoryChangesParams struct {
	// Comma-separated change kinds to include (all by default).
	Kinds []InventoryChangeKind
	// Only changes in these sections (repeat the parameter for several,
	// e.g. the GA sections sold under one floor outline). Names are
	// compared after the seat-map normalisation (letters and digits,
	// lower-cased), so the artifact's section name matches the
	// manifest's (`PIT B` = `PITB`).
	Section []string
	Market  *InventoryMarket
	// Only changes at or after this instant.
	Since *time.Time
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListInventoryChangesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if len(p.Kinds) > 0 {
		q.Set("kinds", joinQuery(p.Kinds, func(v InventoryChangeKind) string { return string(v) }))
	}
	for _, v := range p.Section {
		q.Add("section", v)
	}
	if p.Market != nil {
		q.Set("market", string(*p.Market))
	}
	if p.Since != nil {
		q.Set("since", (*p.Since).Format(time.RFC3339Nano))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListInventoryChanges calls GET /v1/events/{id}/inventory/changes: Recent inventory changes.
//
// Changes between consecutive accepted Ticketmaster primary inventory polls
// of the event, newest first: seats and GA stock that appeared (`restock`,
// `ga_restock`), disappeared (`drop`/`ga_drop`, `pending` while inside the
// cart-hold window), came back from a cart hold (`hold_released`,
// `ga_hold_released`) or settled into sales (`sale`, `ga_sale`). Each row
// groups the seats of one section, row and price level that changed
// together. TM resale activity is listing history
// (`/v1/events/{id}/listings/changes`, source `tmresale`). Pass
// `nextCursor` as `cursor` (with the same filters) for the next page.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListInventoryChanges(ctx context.Context, id int64, params *ListInventoryChangesParams, opts ...RequestOption) (*InventoryChangePage, error) {
	var out InventoryChangePage
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/inventory/changes", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListInventoryChangesPages iterates the pages of ListInventoryChanges, following nextCursor as cursor until the last page.
func (c *Client) ListInventoryChangesPages(ctx context.Context, id int64, params *ListInventoryChangesParams, opts ...RequestOption) iter.Seq2[*InventoryChangePage, error] {
	var p ListInventoryChangesParams
	if params != nil {
		p = *params
	}
	return pages(func() (*InventoryChangePage, error) { return c.ListInventoryChanges(ctx, id, &p, opts...) }, func(page *InventoryChangePage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListInventoryChangesAll iterates every item of ListInventoryChanges across pages.
func (c *Client) ListInventoryChangesAll(ctx context.Context, id int64, params *ListInventoryChangesParams, opts ...RequestOption) iter.Seq2[InventoryChange, error] {
	return items(c.ListInventoryChangesPages(ctx, id, params, opts...), func(page *InventoryChangePage) []InventoryChange { return page.Items })
}

// ListEventPriceLevelsParams holds the query parameters of ListEventPriceLevels.
type ListEventPriceLevelsParams struct {
	// Only this price key.
	PriceKey *int32
}

func (p *ListEventPriceLevelsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.PriceKey != nil {
		q.Set("priceKey", strconv.FormatInt(int64(*p.PriceKey), 10))
	}
	return q
}

// ListEventPriceLevels calls GET /v1/events/{id}/price-levels: Primary price-level history.
//
// Every version of the event's Ticketmaster primary price levels (one per
// AVSC price key): a level whose price, fees or offer changes, or that
// disappears, ends (`validTo`) and its new version starts. Ordered by price
// key, then oldest version first; `validTo` null marks a level in use now.
// Versions that ended before the plan's chart history are left out. Needs a
// plan with map prices.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventPriceLevels(ctx context.Context, id int64, params *ListEventPriceLevelsParams, opts ...RequestOption) (*PriceLevelHistory, error) {
	var out PriceLevelHistory
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/price-levels", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventOffers calls GET /v1/events/{id}/offers: Ticketmaster ticket types.
//
// The offers (ticket types) of the event's Ticketmaster primary price
// levels, each with its levels in use now, their open stock and price
// range: the values `offer` of `GET /v1/events/{id}/series` takes. Offers
// whose levels ended before the plan's chart history are left out; events
// without Ticketmaster price levels answer an empty list. Plans without map
// prices get no prices (`currency`, `minMinor` and `maxMinor` null).
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventOffers(ctx context.Context, id int64, opts ...RequestOption) (*TicketOfferList, error) {
	var out TicketOfferList
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/offers", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventTicketTypes calls GET /v1/events/{id}/ticket-types: Primary ticket types and their changes.
//
// The ticket types the event's primary source sells (Dice: general
// admission at one price per type, without remaining quantities) and the
// latest changes to them: types added or removed, status, price tier and
// all-in price, waiting list. `items` lists the current types in the
// source's order, then removed ones (most recently removed first);
// `changes` are newest first, at most 200. Events whose primary seller
// publishes no ticket types (Ticketmaster: see the price levels) answer
// empty lists. Changes, and types removed, before the plan's chart history
// are left out.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventTicketTypes(ctx context.Context, id int64, opts ...RequestOption) (*TicketTypeList, error) {
	var out TicketTypeList
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/ticket-types", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventListingsParams holds the query parameters of ListEventListings.
type ListEventListingsParams struct {
	// Only this secondary marketplace (`stubhub`); every linked marketplace by default.
	Source *string
	// Only listings of this market (`tm_resale`: Ticketmaster resale; `secondary`: the marketplaces).
	Market *ListEventListingsMarket
	// Only these sections (repeat the parameter for several). A value matches
	// the marketplace's section label after normalisation (letters and digits,
	// lower-cased: `Floor 2` = `FLOOR-2`) or the canonical `sectionKey` exactly.
	Section []string
	// Only listings offering at least this many tickets.
	MinQty *int
	// Only listings whose all-in price per ticket is at most this (minor units of the listing's currency).
	MaxPrice *int64
	// All-in price per ticket, ascending (`price`) or descending (`-price`).
	Sort *ListEventListingsSort
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListEventListingsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	if p.Market != nil {
		q.Set("market", string(*p.Market))
	}
	for _, v := range p.Section {
		q.Add("section", v)
	}
	if p.MinQty != nil {
		q.Set("minQty", strconv.Itoa(*p.MinQty))
	}
	if p.MaxPrice != nil {
		q.Set("maxPrice", strconv.FormatInt(*p.MaxPrice, 10))
	}
	if p.Sort != nil {
		q.Set("sort", string(*p.Sort))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListEventListings calls GET /v1/events/{id}/listings: Current marketplace listings.
//
// The event's current listings on Ticketmaster resale (market `tm_resale`,
// source `tmresale`) and its linked secondary marketplaces, as of each
// source's latest poll (`asOf`), cheapest all-in price per ticket first.
// Marketplace listings need a plan with secondary marketplaces; TM resale
// listings need market data. Pass `nextCursor` as `cursor` (with the same
// filters and sort) for the next page. Live updates: stream topic
// `event.{id}.listings`.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventListings(ctx context.Context, id int64, params *ListEventListingsParams, opts ...RequestOption) (*ListingPage, error) {
	var out ListingPage
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/listings", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventListingsPages iterates the pages of ListEventListings, following nextCursor as cursor until the last page.
func (c *Client) ListEventListingsPages(ctx context.Context, id int64, params *ListEventListingsParams, opts ...RequestOption) iter.Seq2[*ListingPage, error] {
	var p ListEventListingsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*ListingPage, error) { return c.ListEventListings(ctx, id, &p, opts...) }, func(page *ListingPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListEventListingsAll iterates every item of ListEventListings across pages.
func (c *Client) ListEventListingsAll(ctx context.Context, id int64, params *ListEventListingsParams, opts ...RequestOption) iter.Seq2[Listing, error] {
	return items(c.ListEventListingsPages(ctx, id, params, opts...), func(page *ListingPage) []Listing { return page.Items })
}

// ListEventListingChangesParams holds the query parameters of ListEventListingChanges.
type ListEventListingChangesParams struct {
	// Only this secondary marketplace (`stubhub`); every linked marketplace by default.
	Source *string
	// Comma-separated change kinds to include (all by default).
	Kinds []ListingChangeKind
	// Only these sections (repeat the parameter for several). A value matches
	// the marketplace's section label after normalisation (letters and digits,
	// lower-cased: `Floor 2` = `FLOOR-2`) or the canonical `sectionKey` exactly.
	Section []string
	// Only changes at or after this instant.
	Since *time.Time
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListEventListingChangesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	if len(p.Kinds) > 0 {
		q.Set("kinds", joinQuery(p.Kinds, func(v ListingChangeKind) string { return string(v) }))
	}
	for _, v := range p.Section {
		q.Add("section", v)
	}
	if p.Since != nil {
		q.Set("since", (*p.Since).Format(time.RFC3339Nano))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListEventListingChanges calls GET /v1/events/{id}/listings/changes: Recent marketplace listing changes.
//
// Changes between consecutive polls of the event's linked secondary
// marketplaces, newest first: listings that appeared (`added`), changed
// price or quantity, or disappeared (`removed`). Prices are all-in per
// ticket. Needs a plan with secondary marketplaces; history reaches back
// as far as the plan's sales history. Pass `nextCursor` as `cursor`
// (with the same filters) for the next page.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventListingChanges(ctx context.Context, id int64, params *ListEventListingChangesParams, opts ...RequestOption) (*ListingChangePage, error) {
	var out ListingChangePage
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/listings/changes", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventListingChangesPages iterates the pages of ListEventListingChanges, following nextCursor as cursor until the last page.
func (c *Client) ListEventListingChangesPages(ctx context.Context, id int64, params *ListEventListingChangesParams, opts ...RequestOption) iter.Seq2[*ListingChangePage, error] {
	var p ListEventListingChangesParams
	if params != nil {
		p = *params
	}
	return pages(func() (*ListingChangePage, error) { return c.ListEventListingChanges(ctx, id, &p, opts...) }, func(page *ListingChangePage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListEventListingChangesAll iterates every item of ListEventListingChanges across pages.
func (c *Client) ListEventListingChangesAll(ctx context.Context, id int64, params *ListEventListingChangesParams, opts ...RequestOption) iter.Seq2[ListingChange, error] {
	return items(c.ListEventListingChangesPages(ctx, id, params, opts...), func(page *ListingChangePage) []ListingChange { return page.Items })
}

// ListEventListingSectionsParams holds the query parameters of ListEventListingSections.
type ListEventListingSectionsParams struct {
	// Only this secondary marketplace (`stubhub`); every linked marketplace by default.
	Source *string
}

func (p *ListEventListingSectionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	return q
}

// ListEventListingSections calls GET /v1/events/{id}/listings/sections: Current marketplace listings per section.
//
// The event's current marketplace listings aggregated per section, and
// within each section per marketplace (the per-section source compare):
// listings, tickets and the cheapest all-in price per ticket. Sections
// group by canonical `sectionKey`, unmapped labels under their own
// label. Sections are ordered by their cheapest listing (in the minor
// units of each marketplace's currency). Needs a plan with secondary
// marketplaces. Filter the listings of one section with
// `GET /v1/events/{id}/listings?section=<sectionKey or label>`.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventListingSections(ctx context.Context, id int64, params *ListEventListingSectionsParams, opts ...RequestOption) (*ListingSectionList, error) {
	var out ListingSectionList
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/listings/sections", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventMetrics calls GET /v1/events/{id}/metrics: Market metrics of an event.
//
// The event's serving metrics for every scope with market data (KPIs, sparkline,
// movement) that the plan may see: marketplace scopes (`secondary:*`) only with
// secondary marketplaces. Events that started more than 24 hours ago need
// ended-event history (402).
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEventMetrics(ctx context.Context, id int64, opts ...RequestOption) (*EventMetricsList, error) {
	var out EventMetricsList
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/metrics", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventSeriesParams holds the query parameters of GetEventSeries.
type GetEventSeriesParams struct {
	// Chart only this Ticketmaster offer (ticket type), by its price levels' `offerName`; combines with `section`.
	Offer *string
	// Chart these sections as one zone (repeat the parameter for several:
	// `section=101&section=102`): the `key` of the overlay's section
	// layers, the `sectionKey` of the section table. Unknown keys add
	// nothing.
	Section []string
	// Comma-separated metrics (default `floor,median,tickets,sales`).
	Metrics []SeriesMetric
	// Market scope (default `primary`).
	Scope *MetricScope
	Res   *SeriesResolution
	// Range start (default `to` − 72h).
	From *time.Time
	// Range end, exclusive (default now).
	To *time.Time
}

func (p *GetEventSeriesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Offer != nil {
		q.Set("offer", *p.Offer)
	}
	for _, v := range p.Section {
		q.Add("section", v)
	}
	if len(p.Metrics) > 0 {
		q.Set("metrics", joinQuery(p.Metrics, func(v SeriesMetric) string { return string(v) }))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.Res != nil {
		q.Set("res", string(*p.Res))
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	return q
}

// GetEventSeries calls GET /v1/events/{id}/series: Chart series of an event.
//
// Column-oriented buckets of the event's market and sales history from the
// continuous aggregates. Without `res` the resolution follows the range:
// 5m up to 3 days, 1h up to 60 days, 1d beyond. A range above 2000 buckets
// at the requested resolution is a 400.
//
// With `section`, the series is a custom zone: the named sections of the
// `scope` charted together from the section history (the source of
// `GET /v1/events/{id}/sections/{key}/series`). `floor` is the lowest
// zone floor holding at the bucket end, `floorMin` the lowest at any time
// in the bucket (the value held at its start included), `median` the
// priced-ticket-weighted median of the zone sections' medians, `tickets`
// their summed tickets, `sales` the tickets sold in them during the
// bucket and `volume` Σ qty × all-in price of those priced sales, in
// `currency` (the zone's most-listed currency; other currencies convert
// at the day's rate). Range, resolution and plan limits are the same.
// Plans without map prices (`mapPrices`) get a zone's availability only,
// like the section series: `currency` null and every `floor`,
// `floorMin`, `median` and `volume` bucket null.
//
// With `offer` (Ticketmaster events, `scope` primary), the series is one
// ticket type: the primary price levels of that offer (the `name` of
// `GET /v1/events/{id}/offers`), from their seat and GA stock history.
// `floor` is the cheapest of its levels with open stock at the bucket end,
// `floorMin` the cheapest at any time in the bucket, `median` the
// ticket-weighted median of its open levels' all-in prices, `tickets`
// their open stock, `sales` the settled primary sales at those levels
// during the bucket and `volume` Σ qty × all-in price of them. With
// `section` too, only the stock and sales in those sections count (seats
// and GA pools placed through the event's current Ticketmaster manifest).
// Every bucket of the range is returned (null before the event's first
// poll; `floor` and `median` null while nothing is open). `offer` with
// another scope is a 400. Plans without map prices get availability only,
// as for a zone.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetEventSeries(ctx context.Context, id int64, params *GetEventSeriesParams, opts ...RequestOption) (*EventSeries, error) {
	var out EventSeries
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/series", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventSalesParams holds the query parameters of ListEventSales.
type ListEventSalesParams struct {
	// Only this scope (both by default).
	Scope *MetricScope
	// Only sales in these sections (repeat the parameter for several). A
	// value matches the sale's section label after normalisation (letters
	// and digits, lower-cased: `Floor 2` = `FLOOR-2`; Ticketmaster
	// manifest names and marketplace labels) or a marketplace sale's
	// canonical `sectionKey` exactly. No match is an empty page.
	Section []string
	// Comma-separated sale confidences to include (default: every one).
	// `confirmed` sales are reported by a marketplace; Ticketmaster sales
	// are always inferred.
	Confidence []SaleConfidence
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListEventSalesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	for _, v := range p.Section {
		q.Add("section", v)
	}
	if len(p.Confidence) > 0 {
		q.Set("confidence", joinQuery(p.Confidence, func(v SaleConfidence) string { return string(v) }))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListEventSales calls GET /v1/events/{id}/sales: Sales history of an event.
//
// The event's settled sales, newest first (every confidence unless
// `confidence` narrows them; `section` keeps those of some sections).
// Marketplace sales (`secondary:*`) need secondary marketplaces: left
// out without `scope`, 402 when asked for. Pass `nextCursor` as `cursor`
// (same filters) for the next page.
//
// **API key scope:** `sales:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventSales(ctx context.Context, id int64, params *ListEventSalesParams, opts ...RequestOption) (*SalePage, error) {
	var out SalePage
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/sales", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEventSalesPages iterates the pages of ListEventSales, following nextCursor as cursor until the last page.
func (c *Client) ListEventSalesPages(ctx context.Context, id int64, params *ListEventSalesParams, opts ...RequestOption) iter.Seq2[*SalePage, error] {
	var p ListEventSalesParams
	if params != nil {
		p = *params
	}
	return pages(func() (*SalePage, error) { return c.ListEventSales(ctx, id, &p, opts...) }, func(page *SalePage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListEventSalesAll iterates every item of ListEventSales across pages.
func (c *Client) ListEventSalesAll(ctx context.Context, id int64, params *ListEventSalesParams, opts ...RequestOption) iter.Seq2[Sale, error] {
	return items(c.ListEventSalesPages(ctx, id, params, opts...), func(page *SalePage) []Sale { return page.Items })
}

// ListEventSalesSectionsParams holds the query parameters of ListEventSalesSections.
type ListEventSalesSectionsParams struct {
	// Only this scope (every scope by default).
	Scope *MetricScope
	// Comma-separated sale confidences to include (default: every one).
	// `confirmed` sales are reported by a marketplace; Ticketmaster sales
	// are always inferred.
	Confidence []SaleConfidence
}

func (p *ListEventSalesSectionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if len(p.Confidence) > 0 {
		q.Set("confidence", joinQuery(p.Confidence, func(v SaleConfidence) string { return string(v) }))
	}
	return q
}

// ListEventSalesSections calls GET /v1/events/{id}/sales/sections: Sales of an event per section.
//
// The event's settled sales (the same rows as `GET /v1/events/{id}/sales`
// with the same filters) aggregated per section and currency, most
// tickets first. Ticketmaster primary sales group under the manifest's
// section name, TM resale and marketplace sales under their canonical
// `sectionKey` (their own label when unmapped). Unpriced sales (GA units without a price level)
// count in `sales` and `tickets` and join the section's priced currency.
// Without secondary marketplaces, marketplace sales are left out.
//
// **API key scope:** `sales:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListEventSalesSections(ctx context.Context, id int64, params *ListEventSalesSectionsParams, opts ...RequestOption) (*SectionSalesList, error) {
	var out SectionSalesList
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/sales/sections", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSalesFeedParams holds the query parameters of ListSalesFeed.
type ListSalesFeedParams struct {
	// Region code (`GET /v1/regions`); only venues in one of its countries. With `country`, both apply.
	Region *string
	// Venue country (ISO 3166-1 alpha-2).
	Country *string
	// Event segment (case-insensitive), e.g. `Music`.
	Segment *string
	// Minimum all-in price per ticket, minor units of the sale's currency.
	MinPrice *int64
	// Only this scope (both by default).
	Scope *MetricScope
	// Comma-separated sale confidences to include (default: every one).
	// `confirmed` sales are reported by a marketplace; Ticketmaster sales
	// are always inferred.
	Confidence []SaleConfidence
	// Only sales of this source: `ticketmaster` (primary and TM resale) or a
	// marketplace id (`stubhub`, its `secondary:<source>` sales). Combines with
	// `scope` (a scope of another source matches nothing).
	Source *string
	// Only sales of the events in this saved list of the active organization.
	ListID *string
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListSalesFeedParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Region != nil {
		q.Set("region", *p.Region)
	}
	if p.Country != nil {
		q.Set("country", *p.Country)
	}
	if p.Segment != nil {
		q.Set("segment", *p.Segment)
	}
	if p.MinPrice != nil {
		q.Set("minPrice", strconv.FormatInt(*p.MinPrice, 10))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if len(p.Confidence) > 0 {
		q.Set("confidence", joinQuery(p.Confidence, func(v SaleConfidence) string { return string(v) }))
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	if p.ListID != nil {
		q.Set("listId", *p.ListID)
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListSalesFeed calls GET /v1/sales/feed: Sales feed across events.
//
// Settled sales of every event, newest first, each with its event summary:
// Ticketmaster primary and resale sales and marketplace sales (scope
// `secondary:<source>`; `scope=secondary:best` selects every marketplace).
// Marketplace scopes need secondary marketplaces: left out without
// `scope`, 402 when asked for. Every confidence is included unless
// `confidence` narrows it. Live updates arrive on the `sales` stream
// topic (https://status302.com/docs/stream/#topics).
//
// **API key scope:** `sales:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListSalesFeed(ctx context.Context, params *ListSalesFeedParams, opts ...RequestOption) (*SalePage, error) {
	var out SalePage
	if err := c.do(ctx, http.MethodGet, "/v1/sales/feed", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSalesFeedPages iterates the pages of ListSalesFeed, following nextCursor as cursor until the last page.
func (c *Client) ListSalesFeedPages(ctx context.Context, params *ListSalesFeedParams, opts ...RequestOption) iter.Seq2[*SalePage, error] {
	var p ListSalesFeedParams
	if params != nil {
		p = *params
	}
	return pages(func() (*SalePage, error) { return c.ListSalesFeed(ctx, &p, opts...) }, func(page *SalePage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListSalesFeedAll iterates every item of ListSalesFeed across pages.
func (c *Client) ListSalesFeedAll(ctx context.Context, params *ListSalesFeedParams, opts ...RequestOption) iter.Seq2[Sale, error] {
	return items(c.ListSalesFeedPages(ctx, params, opts...), func(page *SalePage) []Sale { return page.Items })
}

// GetSalesSummaryParams holds the query parameters of GetSalesSummary.
type GetSalesSummaryParams struct {
	// How far back the totals reach. Up to `24h` the window starts
	// exactly that long before now; `72h` and `7d` start at the UTC
	// hour holding that instant (`from`). A window longer than the
	// plan's sales history is a 402 (`sales_history`).
	Window *SalesSummaryWindow
	// Region code (`GET /v1/regions`); only venues in one of its countries. With `country`, both apply.
	Region *string
	// Venue country (ISO 3166-1 alpha-2).
	Country *string
	// Event segment (case-insensitive), e.g. `Music`.
	Segment *string
	// Minimum all-in price per ticket, minor units of the sale's currency.
	MinPrice *int64
	// Only this scope (both by default).
	Scope *MetricScope
	// Comma-separated sale confidences to include (default: every one).
	// `confirmed` sales are reported by a marketplace; Ticketmaster sales
	// are always inferred.
	Confidence []SaleConfidence
	// Only sales of this source: `ticketmaster` (primary and TM resale) or a
	// marketplace id (`stubhub`, its `secondary:<source>` sales). Combines with
	// `scope` (a scope of another source matches nothing).
	Source *string
	// Only sales of the events in this saved list of the active organization.
	ListID *string
}

func (p *GetSalesSummaryParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Window != nil {
		q.Set("window", string(*p.Window))
	}
	if p.Region != nil {
		q.Set("region", *p.Region)
	}
	if p.Country != nil {
		q.Set("country", *p.Country)
	}
	if p.Segment != nil {
		q.Set("segment", *p.Segment)
	}
	if p.MinPrice != nil {
		q.Set("minPrice", strconv.FormatInt(*p.MinPrice, 10))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if len(p.Confidence) > 0 {
		q.Set("confidence", joinQuery(p.Confidence, func(v SaleConfidence) string { return string(v) }))
	}
	if p.Source != nil {
		q.Set("source", *p.Source)
	}
	if p.ListID != nil {
		q.Set("listId", *p.ListID)
	}
	return q
}

// GetSalesSummary calls GET /v1/sales/summary: Sales totals of a recent window.
//
// Settled sales of the last 24 hours (or `window`) across events, per
// market scope (every scope the plan may see listed, zeros included):
// sales, tickets and volume per currency. Marketplace scopes need
// secondary marketplaces: left out without `scope`, 402 when asked for.
// Filters and the confidence default (every confidence) match
// `GET /v1/sales/feed`.
//
// **API key scope:** `sales:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetSalesSummary(ctx context.Context, params *GetSalesSummaryParams, opts ...RequestOption) (*SalesSummary, error) {
	var out SalesSummary
	if err := c.do(ctx, http.MethodGet, "/v1/sales/summary", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMoversParams holds the query parameters of ListMovers.
type ListMoversParams struct {
	Window    *ListMoversWindow
	Direction *ListMoversDirection
	// Market scope (default `primary`).
	Scope *MetricScope
	// Only this cohort, e.g. `music/3-7d`.
	Cohort *string
	Limit  *int
}

func (p *ListMoversParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Window != nil {
		q.Set("window", string(*p.Window))
	}
	if p.Direction != nil {
		q.Set("direction", string(*p.Direction))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.Cohort != nil {
		q.Set("cohort", *p.Cohort)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListMovers calls GET /v1/movers: Top movers.
//
// Upcoming events with metrics fresher than 24h, ranked by how much they moved
// within their cohort: `24h` by the movement score, `72h`/`7d` by |z| of the
// window's floor change within the cohort. `direction` keeps rising (`up`)
// or falling (`down`) floors. Marketplace scopes need secondary marketplaces (402).
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListMovers(ctx context.Context, params *ListMoversParams, opts ...RequestOption) (*MoverPage, error) {
	var out MoverPage
	if err := c.do(ctx, http.MethodGet, "/v1/movers", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCompareSeriesParams holds the query parameters of GetCompareSeries.
type GetCompareSeriesParams struct {
	// Comma-separated event ids (1–8).
	// Required.
	Events []int64
	// The metric (default `floor`).
	Metric *SeriesMetric
	Align  *GetCompareSeriesAlign
	// Market scope (default `primary`).
	Scope      *MetricScope
	Res        *SeriesResolution
	From       *time.Time
	To         *time.Time
	DaysBefore *int
}

func (p *GetCompareSeriesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if len(p.Events) > 0 {
		q.Set("events", joinQuery(p.Events, func(v int64) string { return strconv.FormatInt(v, 10) }))
	}
	if p.Metric != nil {
		q.Set("metric", string(*p.Metric))
	}
	if p.Align != nil {
		q.Set("align", string(*p.Align))
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.Res != nil {
		q.Set("res", string(*p.Res))
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	if p.DaysBefore != nil {
		q.Set("daysBefore", strconv.Itoa(*p.DaysBefore))
	}
	return q
}

// GetCompareSeries calls GET /v1/compare/series: One metric across several events.
//
// `align=calendar` returns every event over the same range (`from`/`to`,
// default the last 30 days); `align=days_to_event` returns each event over
// the `daysBefore` days before its own start (up to now for upcoming
// shows), with `daysToEvent` as the shared axis.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetCompareSeries(ctx context.Context, params *GetCompareSeriesParams, opts ...RequestOption) (*CompareSeriesSet, error) {
	var out CompareSeriesSet
	if err := c.do(ctx, http.MethodGet, "/v1/compare/series", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchCatalogParams holds the query parameters of SearchCatalog.
type SearchCatalogParams struct {
	// Required.
	Q string
	// Hits per group.
	Limit *int
}

func (p *SearchCatalogParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	q.Set("q", p.Q)
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// SearchCatalog calls GET /v1/search: Typeahead across events, attractions and venues.
//
// Name search for the ⌘K palette. The query and the names are compared
// folded (lower case, accents and punctuation removed, `&` = `and`), by
// substring and pg_trgm word similarity; queries shorter than three
// characters match name and word prefixes only. Each group is ranked:
// exact name, then name prefix, then word prefix, then substring/trigram
// similarity. Events list upcoming ones (not started, or undated) before
// past ones, then by match and soonest start; attractions and venues
// with upcoming events come first, by upcoming event count within the
// exact/prefix tiers and by similarity within the trigram tier.
//
// A `q` that is an event ID or event page URL also finds the events it
// names, listed first: a Ticketmaster event ID (16 hex digits, any case,
// e.g. `380064AEC95E177C`, or a ticketmaster.* `…/event/{id}` URL) or a
// StubHub/viagogo event number of a linked marketplace listing (e.g.
// `161219435`, viagogo's `E-161219435`, or a stubhub.* / viagogo.*
// `…/event/{n}` or `…/E-{n}` URL). URLs may omit the scheme; their query
// string and fragment are ignored.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) SearchCatalog(ctx context.Context, params *SearchCatalogParams, opts ...RequestOption) (*SearchResults, error) {
	var out SearchResults
	if err := c.do(ctx, http.MethodGet, "/v1/search", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAttraction calls GET /v1/attractions/{id}: Attraction detail.
//
// One attraction (performer, team, show) with its Discovery
// classifications and links, a rollup of its catalog events around now
// and its next upcoming event.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetAttraction(ctx context.Context, id int64, opts ...RequestOption) (*AttractionDetail, error) {
	var out AttractionDetail
	if err := c.do(ctx, http.MethodGet, "/v1/attractions/"+url.PathEscape(strconv.FormatInt(id, 10)), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAttractionEventsParams holds the query parameters of ListAttractionEvents.
type ListAttractionEventsParams struct {
	// Only events starting at or after this instant. Defaults to now, in which case undated (TBA/TBD) events are included too.
	From *time.Time
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListAttractionEventsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListAttractionEvents calls GET /v1/attractions/{id}/events: Events of an attraction.
//
// The attraction's catalog events in start order (undated last), upcoming
// by default, with cursor pagination like `listEvents`.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListAttractionEvents(ctx context.Context, id int64, params *ListAttractionEventsParams, opts ...RequestOption) (*EventPage, error) {
	var out EventPage
	if err := c.do(ctx, http.MethodGet, "/v1/attractions/"+url.PathEscape(strconv.FormatInt(id, 10))+"/events", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAttractionEventsPages iterates the pages of ListAttractionEvents, following nextCursor as cursor until the last page.
func (c *Client) ListAttractionEventsPages(ctx context.Context, id int64, params *ListAttractionEventsParams, opts ...RequestOption) iter.Seq2[*EventPage, error] {
	var p ListAttractionEventsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*EventPage, error) { return c.ListAttractionEvents(ctx, id, &p, opts...) }, func(page *EventPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListAttractionEventsAll iterates every item of ListAttractionEvents across pages.
func (c *Client) ListAttractionEventsAll(ctx context.Context, id int64, params *ListAttractionEventsParams, opts ...RequestOption) iter.Seq2[EventSummary, error] {
	return items(c.ListAttractionEventsPages(ctx, id, params, opts...), func(page *EventPage) []EventSummary { return page.Items })
}

// GetAttractionHistoryParams holds the query parameters of GetAttractionHistory.
type GetAttractionHistoryParams struct {
	// Past events returned, most recent first.
	Limit *int
}

func (p *GetAttractionHistoryParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// GetAttractionHistory calls GET /v1/attractions/{id}/history: Market history of an attraction's past events.
//
// The attraction's events that have started, most recent first, each with
// how its market ended: the resale and secondary spread over the primary
// floor and the primary sell-through from the event's last metrics, and
// the primary depletion over its collected history (first to last primary
// sample up to the start). Use it to judge an upcoming on-sale by how the
// act's earlier shows went. Needs market data and ended-event history;
// `secondarySpread` is null without secondary marketplaces.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetAttractionHistory(ctx context.Context, id int64, params *GetAttractionHistoryParams, opts ...RequestOption) (*AttractionHistory, error) {
	var out AttractionHistory
	if err := c.do(ctx, http.MethodGet, "/v1/attractions/"+url.PathEscape(strconv.FormatInt(id, 10))+"/history", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVenue calls GET /v1/venues/{id}: Venue detail.
//
// One venue with its address and location, a rollup of its catalog events
// around now, its next upcoming event, and the Ticketmaster venue
// configurations (layouts) of its events' compiled seat maps.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetVenue(ctx context.Context, id int64, opts ...RequestOption) (*VenueDetail, error) {
	var out VenueDetail
	if err := c.do(ctx, http.MethodGet, "/v1/venues/"+url.PathEscape(strconv.FormatInt(id, 10)), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListVenueEventsParams holds the query parameters of ListVenueEvents.
type ListVenueEventsParams struct {
	// Only events starting at or after this instant. Defaults to now, in which case undated (TBA/TBD) events are included too.
	From *time.Time
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListVenueEventsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListVenueEvents calls GET /v1/venues/{id}/events: Events at a venue.
//
// The venue's catalog events in start order (undated last), upcoming by
// default, with cursor pagination like `listEvents`.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListVenueEvents(ctx context.Context, id int64, params *ListVenueEventsParams, opts ...RequestOption) (*EventPage, error) {
	var out EventPage
	if err := c.do(ctx, http.MethodGet, "/v1/venues/"+url.PathEscape(strconv.FormatInt(id, 10))+"/events", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListVenueEventsPages iterates the pages of ListVenueEvents, following nextCursor as cursor until the last page.
func (c *Client) ListVenueEventsPages(ctx context.Context, id int64, params *ListVenueEventsParams, opts ...RequestOption) iter.Seq2[*EventPage, error] {
	var p ListVenueEventsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*EventPage, error) { return c.ListVenueEvents(ctx, id, &p, opts...) }, func(page *EventPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListVenueEventsAll iterates every item of ListVenueEvents across pages.
func (c *Client) ListVenueEventsAll(ctx context.Context, id int64, params *ListVenueEventsParams, opts ...RequestOption) iter.Seq2[EventSummary, error] {
	return items(c.ListVenueEventsPages(ctx, id, params, opts...), func(page *EventPage) []EventSummary { return page.Items })
}

// ListCalendarDaysParams holds the query parameters of ListCalendarDays.
type ListCalendarDaysParams struct {
	// First local day of the range (inclusive) in `tz`.
	// Required.
	From string
	// Local day after the range (exclusive) in `tz`; must be after `from`.
	// Required.
	To string
	// IANA time zone of the viewer; days and `date` are local to it. Unknown zones are rejected.
	TZ *string
	// Comma-separated item kinds to include (all by default).
	Kinds []CalendarKind
	// Venue country (ISO 3166-1 alpha-2).
	Country *string
	// Region code (`GET /v1/regions`); only venues in one of its countries. With `country`, both apply.
	Region *string
	// Only events of this primary seller (`CalendarEvent.primarySource`): `tm` for events hosted on Ticketmaster, `dice` for Dice's own events.
	PrimarySource *PrimarySource
	// Only events hosted on Ticketmaster (with a `tmEventId`). Same as `primarySource=tm`, which replaces it; 400 together with another `primarySource`.
	TmOnly *bool
	// Event name search (substring or trigram word similarity on the normalized name).
	Q *string
	// Only events in any of the active organization's saved lists.
	MyLists *bool
	// Items returned per day, most prominent first (a month cell shows a few, the week grid more).
	Top *int
}

func (p *ListCalendarDaysParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	q.Set("from", p.From)
	q.Set("to", p.To)
	if p.TZ != nil {
		q.Set("tz", *p.TZ)
	}
	if len(p.Kinds) > 0 {
		q.Set("kinds", joinQuery(p.Kinds, func(v CalendarKind) string { return string(v) }))
	}
	if p.Country != nil {
		q.Set("country", *p.Country)
	}
	if p.Region != nil {
		q.Set("region", *p.Region)
	}
	if p.PrimarySource != nil {
		q.Set("primarySource", string(*p.PrimarySource))
	}
	if p.TmOnly != nil {
		q.Set("tmOnly", strconv.FormatBool(*p.TmOnly))
	}
	if p.Q != nil {
		q.Set("q", *p.Q)
	}
	if p.MyLists != nil {
		q.Set("myLists", strconv.FormatBool(*p.MyLists))
	}
	if p.Top != nil {
		q.Set("top", strconv.Itoa(*p.Top))
	}
	return q
}

// ListCalendarDays calls GET /v1/calendar/days: Calendar items counted per day.
//
// One row per local day of `[from, to)` in `tz` that has at least one item
// (days without items are omitted), with the item count per kind and the
// `top` most prominent items of the day (see `CalendarItem.prominence`;
// an event's first item of the day ranks ahead of its later ones).
// Items are public on-sale starts (`onsale`), presale starts (`presale`,
// presales of one event that start together are one item) and show
// starts (`event`). Items without a known instant are excluded: TBA/TBD
// dates, shows whose time is TBA, public on-sales flagged TBD and presales
// without a start. The range is at most 62 days (a 6-week month grid).
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListCalendarDays(ctx context.Context, params *ListCalendarDaysParams, opts ...RequestOption) (*CalendarDays, error) {
	var out CalendarDays
	if err := c.do(ctx, http.MethodGet, "/v1/calendar/days", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListCalendarItemsParams holds the query parameters of ListCalendarItems.
type ListCalendarItemsParams struct {
	// First local day of the range (inclusive) in `tz`.
	// Required.
	From string
	// Local day after the range (exclusive) in `tz`; must be after `from`.
	// Required.
	To string
	// IANA time zone of the viewer; days and `date` are local to it. Unknown zones are rejected.
	TZ *string
	// Comma-separated item kinds to include (all by default).
	Kinds []CalendarKind
	// Venue country (ISO 3166-1 alpha-2).
	Country *string
	// Region code (`GET /v1/regions`); only venues in one of its countries. With `country`, both apply.
	Region *string
	// Only events of this primary seller (`CalendarEvent.primarySource`): `tm` for events hosted on Ticketmaster, `dice` for Dice's own events.
	PrimarySource *PrimarySource
	// Only events hosted on Ticketmaster (with a `tmEventId`). Same as `primarySource=tm`, which replaces it; 400 together with another `primarySource`.
	TmOnly *bool
	// Event name search (substring or trigram word similarity on the normalized name).
	Q *string
	// Only events in any of the active organization's saved lists.
	MyLists *bool
	Sort    *ListCalendarItemsSort
	// Opaque `nextCursor` of the previous page.
	Cursor *string
	Limit  *int
}

func (p *ListCalendarItemsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	q.Set("from", p.From)
	q.Set("to", p.To)
	if p.TZ != nil {
		q.Set("tz", *p.TZ)
	}
	if len(p.Kinds) > 0 {
		q.Set("kinds", joinQuery(p.Kinds, func(v CalendarKind) string { return string(v) }))
	}
	if p.Country != nil {
		q.Set("country", *p.Country)
	}
	if p.Region != nil {
		q.Set("region", *p.Region)
	}
	if p.PrimarySource != nil {
		q.Set("primarySource", string(*p.PrimarySource))
	}
	if p.TmOnly != nil {
		q.Set("tmOnly", strconv.FormatBool(*p.TmOnly))
	}
	if p.Q != nil {
		q.Set("q", *p.Q)
	}
	if p.MyLists != nil {
		q.Set("myLists", strconv.FormatBool(*p.MyLists))
	}
	if p.Sort != nil {
		q.Set("sort", string(*p.Sort))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListCalendarItems calls GET /v1/calendar/items: Calendar items.
//
// On-sale, presale and show items starting in the local days `[from, to)`
// of `tz` (at most 366 days), with cursor pagination. `sort=time` (the
// default) is chronological; `sort=prominence` puts the most prominent
// first. Pass `nextCursor` as `cursor` with the same parameters for the
// next page. The same exclusions as `/v1/calendar/days` apply.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListCalendarItems(ctx context.Context, params *ListCalendarItemsParams, opts ...RequestOption) (*CalendarItemPage, error) {
	var out CalendarItemPage
	if err := c.do(ctx, http.MethodGet, "/v1/calendar/items", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListCalendarItemsPages iterates the pages of ListCalendarItems, following nextCursor as cursor until the last page.
func (c *Client) ListCalendarItemsPages(ctx context.Context, params *ListCalendarItemsParams, opts ...RequestOption) iter.Seq2[*CalendarItemPage, error] {
	var p ListCalendarItemsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*CalendarItemPage, error) { return c.ListCalendarItems(ctx, &p, opts...) }, func(page *CalendarItemPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListCalendarItemsAll iterates every item of ListCalendarItems across pages.
func (c *Client) ListCalendarItemsAll(ctx context.Context, params *ListCalendarItemsParams, opts ...RequestOption) iter.Seq2[CalendarItem, error] {
	return items(c.ListCalendarItemsPages(ctx, params, opts...), func(page *CalendarItemPage) []CalendarItem { return page.Items })
}

// GetSeatmap calls GET /v1/seatmaps/{hash}: Seat map artifact.
//
// The compiled `.tmap` seat map artifact,
// addressed by the SHA-256 of its uncompressed bytes and therefore
// immutable. Sent brotli-compressed (`Content-Encoding: br`) when the
// request accepts `br`, uncompressed otherwise. The ETag is the hash.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetSeatmap(ctx context.Context, hash string, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := c.do(ctx, http.MethodGet, "/v1/seatmaps/"+url.PathEscape(hash), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}

// ListPredictionsParams holds the query parameters of ListPredictions.
type ListPredictionsParams struct {
	// The call family (`tm`, the default, or `dice`).
	Category *PredictionCategory
	Status   *PredictionStatusFilter
	// The `nextCursor` of the previous page (with the same `status`).
	Cursor *string
	Limit  *int
}

func (p *ListPredictionsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Category != nil {
		q.Set("category", string(*p.Category))
	}
	if p.Status != nil {
		q.Set("status", string(*p.Status))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListPredictions calls GET /v1/predictions: Scalping predictions.
//
// Calls on upcoming events whose primary tickets are expected to resell
// above the price paid, one category at a time (`category`: `tm`, the
// default, for Ticketmaster events; `dice` for DICE events resold on a
// marketplace): the event, the entry (the primary floor when the call
// was made), the risk (the chance the call is wrong, 0–100) with its
// reasons, the latest resale floor and markup, and once graded the
// exit, the exit markup and the outcome. `status=open` (the default)
// lists open calls lowest risk first; `graded` (hits and misses),
// `hit`, `miss` and `void` list graded calls most recently graded
// first. Money is in minor units of `currency` (see the money guide).
//
// Every plan includes predictions; a category answers 403
// `https://status302.com/problems/predictions-unpublished` until it
// is published, and `GET /v1/predictions/availability` tells when. API
// keys and signed-in members see the same calls. The calls are model
// output, not financial advice: do your own research.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListPredictions(ctx context.Context, params *ListPredictionsParams, opts ...RequestOption) (*PredictionCallPage, error) {
	var out PredictionCallPage
	if err := c.do(ctx, http.MethodGet, "/v1/predictions", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPredictionsPages iterates the pages of ListPredictions, following nextCursor as cursor until the last page.
func (c *Client) ListPredictionsPages(ctx context.Context, params *ListPredictionsParams, opts ...RequestOption) iter.Seq2[*PredictionCallPage, error] {
	var p ListPredictionsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*PredictionCallPage, error) { return c.ListPredictions(ctx, &p, opts...) }, func(page *PredictionCallPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListPredictionsAll iterates every item of ListPredictions across pages.
func (c *Client) ListPredictionsAll(ctx context.Context, params *ListPredictionsParams, opts ...RequestOption) iter.Seq2[PredictionCall, error] {
	return items(c.ListPredictionsPages(ctx, params, opts...), func(page *PredictionCallPage) []PredictionCall { return page.Items })
}

// GetPredictionTrackRecordParams holds the query parameters of GetPredictionTrackRecord.
type GetPredictionTrackRecordParams struct {
	// The call family (`tm`, the default, or `dice`).
	Category *PredictionCategory
}

func (p *GetPredictionTrackRecordParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Category != nil {
		q.Set("category", string(*p.Category))
	}
	return q
}

// GetPredictionTrackRecord calls GET /v1/predictions/track-record: Prediction track record.
//
// How one category's calls turned out, overall and per model version,
// all time and over the last 30 and 90 days (by grading time): hits,
// misses, voids, hit rate, mean exit markup, Brier score of the risk at
// the call, and the calibration per 10-point risk band (the miss rate
// the risk predicted against the one observed). Recomputed at most
// every minute. Same access as `GET /v1/predictions` for the category.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetPredictionTrackRecord(ctx context.Context, params *GetPredictionTrackRecordParams, opts ...RequestOption) (*PredictionTrackRecord, error) {
	var out PredictionTrackRecord
	if err := c.do(ctx, http.MethodGet, "/v1/predictions/track-record", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPredictionAvailability calls GET /v1/predictions/availability: Whether the caller sees the predictions.
//
// Customers on any plan see a category of predictions once it is
// published: staff switched it on for the org and its current model has
// a track record of at least `minGraded` graded calls. Each category
// (`tm`, `dice`) is published on its own; the top-level fields are
// `tm`'s. Any member or key may ask, so a client can tell a wait
// (`reason` `unpublished`) before calling `GET /v1/predictions`.
//
// **API key scope:** `market:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetPredictionAvailability(ctx context.Context, opts ...RequestOption) (*PredictionAvailability, error) {
	var out PredictionAvailability
	if err := c.do(ctx, http.MethodGet, "/v1/predictions/availability", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMustCopsParams holds the query parameters of ListMustCops.
type ListMustCopsParams struct {
	Status *MustCopStatusFilter
	// The `nextCursor` of the previous page (with the same `status`).
	Cursor *string
	Limit  *int
}

func (p *ListMustCopsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Status != nil {
		q.Set("status", string(*p.Status))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListMustCops calls GET /v1/must-cop: Must cops.
//
// Upcoming public on-sales (Ticketmaster or DICE, within 60 days) of acts
// whose latest known on-sales sold out instantly: at least 2 of the act's
// last 3 (the last two years, both sellers; the act is the headliner)
// sold out within two hours of their public on-sale. Each pick carries
// the on-sales it rests on and is graded from its own on-sale: `hit`
// when it sold out within two hours too. `status=open` (the default)
// lists open picks soonest on-sale first; `graded` (hits and misses),
// `hit`, `miss` and `void` list graded picks most recently graded
// first. `record` counts every graded pick. Refreshed every 15 minutes.
//
// Every plan includes them once they are published to the org (403
// `https://status302.com/problems/must-cop-unpublished` until staff
// do); staff always see them. API keys and signed-in members see the same
// picks (keys get no staff bypass). Not financial advice: do your own
// research.
//
// **API key scope:** `market:read`.
func (c *Client) ListMustCops(ctx context.Context, params *ListMustCopsParams, opts ...RequestOption) (*MustCopPage, error) {
	var out MustCopPage
	if err := c.do(ctx, http.MethodGet, "/v1/must-cop", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMustCopsPages iterates the pages of ListMustCops, following nextCursor as cursor until the last page.
func (c *Client) ListMustCopsPages(ctx context.Context, params *ListMustCopsParams, opts ...RequestOption) iter.Seq2[*MustCopPage, error] {
	var p ListMustCopsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*MustCopPage, error) { return c.ListMustCops(ctx, &p, opts...) }, func(page *MustCopPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListMustCopsAll iterates every item of ListMustCops across pages.
func (c *Client) ListMustCopsAll(ctx context.Context, params *ListMustCopsParams, opts ...RequestOption) iter.Seq2[MustCop, error] {
	return items(c.ListMustCopsPages(ctx, params, opts...), func(page *MustCopPage) []MustCop { return page.Items })
}

// ListSavedLists calls GET /v1/lists: The active org's saved lists.
//
// Every list of the active org, by name, with item counts. Any member reads them.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListSavedLists(ctx context.Context, opts ...RequestOption) (*SavedListList, error) {
	var out SavedListList
	if err := c.do(ctx, http.MethodGet, "/v1/lists", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateSavedList calls POST /v1/lists: Create a list.
//
// Creates a list in the active org, optionally with its first events (at
// most 500). The org's plan limits the number of lists and of events per
// list (402). Recorded in the audit log.
//
// **API key scope:** `lists:write`.
func (c *Client) CreateSavedList(ctx context.Context, body *CreateSavedListRequest, opts ...RequestOption) (*SavedList, error) {
	var out SavedList
	if err := c.do(ctx, http.MethodPost, "/v1/lists", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSavedList calls GET /v1/lists/{listId}: One list with its saved view.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetSavedList(ctx context.Context, listID string, opts ...RequestOption) (*SavedList, error) {
	var out SavedList
	if err := c.do(ctx, http.MethodGet, "/v1/lists/"+url.PathEscape(listID), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSavedList calls PATCH /v1/lists/{listId}: Rename a list or change its saved view.
//
// Any member. View fields that are present replace the stored ones.
//
// **API key scope:** `lists:write`.
func (c *Client) UpdateSavedList(ctx context.Context, listID string, body *UpdateSavedListRequest, opts ...RequestOption) (*SavedList, error) {
	var out SavedList
	if err := c.do(ctx, http.MethodPatch, "/v1/lists/"+url.PathEscape(listID), nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSavedList calls DELETE /v1/lists/{listId}: Delete a list.
//
// The list's creator or an org admin/owner. Recorded in the audit log.
//
// **API key scope:** `lists:write`.
func (c *Client) DeleteSavedList(ctx context.Context, listID string, opts ...RequestOption) error {
	return c.do(ctx, http.MethodDelete, "/v1/lists/"+url.PathEscape(listID), nil, nil, nil, opts)
}

// ListSavedListItems calls GET /v1/lists/{listId}/items: The events of a list with their user columns.
//
// Oldest first.
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) ListSavedListItems(ctx context.Context, listID string, opts ...RequestOption) (*SavedListItemList, error) {
	var out SavedListItemList
	if err := c.do(ctx, http.MethodGet, "/v1/lists/"+url.PathEscape(listID)+"/items", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddSavedListItems calls POST /v1/lists/{listId}/items: Add events to a list.
//
// Idempotent bulk add (1–500 events): events already in the list are
// skipped. Unknown event ids reject the whole call (400). The plan's
// per-list limit applies to the events new to the list (402).
//
// **API key scope:** `lists:write`.
func (c *Client) AddSavedListItems(ctx context.Context, listID string, body *AddSavedListItemsRequest, opts ...RequestOption) (*AddSavedListItemsResult, error) {
	var out AddSavedListItemsResult
	if err := c.do(ctx, http.MethodPost, "/v1/lists/"+url.PathEscape(listID)+"/items", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSavedListItems calls PATCH /v1/lists/{listId}/items: Edit the user columns of items.
//
// Notes, cost basis (per ticket, minor units) and quantity held of 1–500
// items, all or nothing. A cost basis without `currency` keeps the item's
// currency or takes the event's market currency. `clear` empties
// columns. An event that is not in the list is a 404.
//
// **API key scope:** `lists:write`.
func (c *Client) UpdateSavedListItems(ctx context.Context, listID string, body *UpdateSavedListItemsRequest, opts ...RequestOption) (*SavedListItemList, error) {
	var out SavedListItemList
	if err := c.do(ctx, http.MethodPatch, "/v1/lists/"+url.PathEscape(listID)+"/items", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveSavedListItemsParams holds the query parameters of RemoveSavedListItems.
type RemoveSavedListItemsParams struct {
	// Comma-separated event ids (1–500).
	// Required.
	EventIDs []int64
}

func (p *RemoveSavedListItemsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if len(p.EventIDs) > 0 {
		q.Set("eventIds", joinQuery(p.EventIDs, func(v int64) string { return strconv.FormatInt(v, 10) }))
	}
	return q
}

// RemoveSavedListItems calls DELETE /v1/lists/{listId}/items: Remove events from a list.
//
// **API key scope:** `lists:write`.
func (c *Client) RemoveSavedListItems(ctx context.Context, listID string, params *RemoveSavedListItemsParams, opts ...RequestOption) (*RemoveSavedListItemsResult, error) {
	var out RemoveSavedListItemsResult
	if err := c.do(ctx, http.MethodDelete, "/v1/lists/"+url.PathEscape(listID)+"/items", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExportSavedListCSVParams holds the query parameters of ExportSavedListCSV.
type ExportSavedListCSVParams struct {
	// Market scope (default `primary`).
	Scope *MetricScope
	// Comma-separated column ids in display order (`GridColumn.id`); absent = the saved view's, else the defaults.
	Columns []string
	// Comma-separated sort keys, first wins; a leading `-` sorts descending (`-chg24h,name`). Empty = start date; absent = the saved view's.
	Sort *string
	// Repeatable `column:op[:value]`. Numeric ops `lt`, `lte`, `gt`, `gte`, `eq`, `ne` take numbers in API units
	// (minor units, fractions, unix seconds); text ops `contains`, `notContains`, `eq`, `ne`; any column
	// `empty`, `notEmpty`. Empty entries are ignored, so `filter=` clears the saved filters.
	Filter []string
	// Share of the resale price the seller pays, for the `roi` column (default 0; the web app sends the
	// user's seller fee for the resale market, `GET /v1/me/preferences`).
	SellerFee *float64
}

func (p *ExportSavedListCSVParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if len(p.Columns) > 0 {
		q.Set("columns", joinQuery(p.Columns, func(v string) string { return v }))
	}
	if p.Sort != nil {
		q.Set("sort", *p.Sort)
	}
	for _, v := range p.Filter {
		q.Add("filter", v)
	}
	if p.SellerFee != nil {
		q.Set("sellerFee", strconv.FormatFloat(*p.SellerFee, 'f', -1, 64))
	}
	return q
}

// ExportSavedListCSV calls GET /v1/lists/{listId}/export.csv: The grid as CSV.
//
// The same rows as the grid (same parameters and fallbacks) as RFC 4180
// CSV: event id and name, the columns in order (sparklines left out),
// then the row currency. Money in major units, changes and shares as
// percentages, times in RFC 3339 UTC. Text starting with `= + - @` is
// prefixed with an apostrophe (formula injection). Counts against the
// same in-flight limit as the grid (429, `limit` `concurrent`).
//
// **API key scope:** `catalog:read`.
func (c *Client) ExportSavedListCSV(ctx context.Context, listID string, params *ExportSavedListCSVParams, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := c.do(ctx, http.MethodGet, "/v1/lists/"+url.PathEscape(listID)+"/export.csv", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}

// ListEventSavedLists calls GET /v1/events/{id}/lists: The active org's lists containing an event.
//
// **API key scope:** `catalog:read`.
func (c *Client) ListEventSavedLists(ctx context.Context, id int64, opts ...RequestOption) (*EventSavedLists, error) {
	var out EventSavedLists
	if err := c.do(ctx, http.MethodGet, "/v1/events/"+url.PathEscape(strconv.FormatInt(id, 10))+"/lists", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPortfolioParams holds the query parameters of GetPortfolio.
type GetPortfolioParams struct {
	// Market scope (default `primary`).
	Scope *MetricScope
	// Rows per list item held (`position`, default) or per event (`event`).
	GroupBy *PortfolioGroupBy
}

func (p *GetPortfolioParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.GroupBy != nil {
		q.Set("groupBy", string(*p.GroupBy))
	}
	return q
}

// GetPortfolio calls GET /v1/portfolio: The active org's positions across its lists.
//
// Every list item with a quantity held, valued at the floor of `scope`
// (the grid's position columns: cost = qty × cost basis, value = qty ×
// floor, P&L = qty × (floor − cost basis); no seller fee is deducted).
// An event held in two lists is two positions (`eventPositions` counts
// them); `groupBy=event` folds them into one row per event (`lists`
// names its lists; quantities and money summed, `costBasis` the
// average per ticket). `totals` equal the overview's `holdings`; `exposure` splits
// them by the event's date: `past` (the `past`/`archived` lifecycles —
// tickets still held after the event), `next7d`, `next30d` (8–30 days),
// `later` and `undated`. Positions in upcoming events carry at-risk
// reasons. Any member or key reads it; market data in the scope needs
// the plan's (402). Money is in minor units (cost in `costCurrency`,
// market value and P&L in the market's currency; see the money guide).
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetPortfolio(ctx context.Context, params *GetPortfolioParams, opts ...RequestOption) (*Portfolio, error) {
	var out Portfolio
	if err := c.do(ctx, http.MethodGet, "/v1/portfolio", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExportPortfolioCSVParams holds the query parameters of ExportPortfolioCSV.
type ExportPortfolioCSVParams struct {
	// Market scope (default `primary`).
	Scope *MetricScope
	// Rows per list item held (`position`, default) or per event (`event`).
	GroupBy *PortfolioGroupBy
}

func (p *ExportPortfolioCSVParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.GroupBy != nil {
		q.Set("groupBy", string(*p.GroupBy))
	}
	return q
}

// ExportPortfolioCSV calls GET /v1/portfolio/export.csv: The portfolio's positions as CSV.
//
// The rows of `GET /v1/portfolio` (per position, or per event with
// `groupBy=event`: lists joined by "; ", average cost basis) in its order as RFC 4180 CSV: event,
// list, quantity, cost basis and total (cost currency), floor, value,
// P&L and P&L % (market currency), Δ24h, sales 24h, bucket, at-risk
// reasons and the market sample time. Money in major units, changes as
// percentages, times in RFC 3339 UTC; text starting with `= + - @` is
// prefixed with an apostrophe. Needs CSV export and the plan's export
// row limit, like list exports (402).
//
// **API key scope:** `catalog:read`.
func (c *Client) ExportPortfolioCSV(ctx context.Context, params *ExportPortfolioCSVParams, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := c.do(ctx, http.MethodGet, "/v1/portfolio/export.csv", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPortfolioHistoryParams holds the query parameters of GetPortfolioHistory.
type GetPortfolioHistoryParams struct {
	// Market scope (default `primary`).
	Scope    *MetricScope
	From     *time.Time
	To       *time.Time
	Interval *PortfolioHistoryInterval
}

func (p *GetPortfolioHistoryParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Scope != nil {
		q.Set("scope", string(*p.Scope))
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	if p.Interval != nil {
		q.Set("interval", string(*p.Interval))
	}
	return q
}

// GetPortfolioHistory calls GET /v1/portfolio/history: The portfolio's cost, market value and P&L over time.
//
// The org's held list items summed per bucket of `interval` over
// `[from, to)`, one series per currency, with the meanings of the
// portfolio's `totals`: `cost` in the cost basis's currency, `value`
// (qty × floor) and `pnl` (qty × (floor − cost basis)) in the market's,
// `pnlCost` the cost of the positions in `pnl`. Floors come from the
// hourly or daily market aggregates of `scope`: a bucket's floor is its
// last one, else the last before it (up to 7 days before `from`).
// Assumptions: an item counts from the bucket in which it was added to
// its list, with its current quantity and cost basis (earlier
// quantities and removed items are not recorded); positions in events
// that ended beyond the plan's history keep their cost but no value.
// `from` defaults to 30 days before `to` (default now); `interval`
// defaults to `1h` for spans up to 31 days, `1d` beyond; at most 2000
// buckets. The range is clamped to the plan's series history (402 when
// nothing remains). Market data in the scope needs the plan's (402).
//
// **API key scope:** `catalog:read`.
//
// Conditional: with WithIfNoneMatch(etag) the call fails with ErrNotModified (*NotModifiedError) while the representation still has that ETag.
func (c *Client) GetPortfolioHistory(ctx context.Context, params *GetPortfolioHistoryParams, opts ...RequestOption) (*PortfolioHistory, error) {
	var out PortfolioHistory
	if err := c.do(ctx, http.MethodGet, "/v1/portfolio/history", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAlertRulesParams holds the query parameters of ListAlertRules.
type ListAlertRulesParams struct {
	// Rules older than this id (`nextBefore`).
	Before *int64
	Limit  *int
}

func (p *ListAlertRulesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Before != nil {
		q.Set("before", strconv.FormatInt(*p.Before, 10))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListAlertRules calls GET /v1/alerts/rules: The org's alert rules.
//
// Newest first, with the org's rule count and plan limit.
//
// **API key scope:** `alerts:write`.
func (c *Client) ListAlertRules(ctx context.Context, params *ListAlertRulesParams, opts ...RequestOption) (*AlertRulePage, error) {
	var out AlertRulePage
	if err := c.do(ctx, http.MethodGet, "/v1/alerts/rules", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAlertRulesPages iterates the pages of ListAlertRules, following nextBefore as before until the last page.
func (c *Client) ListAlertRulesPages(ctx context.Context, params *ListAlertRulesParams, opts ...RequestOption) iter.Seq2[*AlertRulePage, error] {
	var p ListAlertRulesParams
	if params != nil {
		p = *params
	}
	return pages(func() (*AlertRulePage, error) { return c.ListAlertRules(ctx, &p, opts...) }, func(page *AlertRulePage) bool {
		if page.NextBefore == nil {
			return false
		}
		next := int64(*page.NextBefore)
		p.Before = new(next)
		return true
	})
}

// ListAlertRulesAll iterates every item of ListAlertRules across pages.
func (c *Client) ListAlertRulesAll(ctx context.Context, params *ListAlertRulesParams, opts ...RequestOption) iter.Seq2[AlertRule, error] {
	return items(c.ListAlertRulesPages(ctx, params, opts...), func(page *AlertRulePage) []AlertRule { return page.Items })
}

// CreateAlertRule calls POST /v1/alerts/rules: Create a rule.
//
// Validates the scope (event, list, attraction or catalog filter) and the
// condition against the plan (rule count, % change metrics, scopes).
// Metric rules are evaluated against current metrics right away
// (events already matching fire once).
//
// **API key scope:** `alerts:write`.
func (c *Client) CreateAlertRule(ctx context.Context, body *AlertRuleInput, opts ...RequestOption) (*AlertRule, error) {
	var out AlertRule
	if err := c.do(ctx, http.MethodPost, "/v1/alerts/rules", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// PreviewAlertRule calls POST /v1/alerts/rules/preview: Preview a draft rule.
//
// Dry-runs a draft rule (the body of `POST /v1/alerts/rules`; its
// `channels` are ignored) so you can tune a threshold before saving
// it. Metric conditions list the events in scope matching now (section
// conditions: one item per matching section of the first 50 upcoming
// events in scope, `detail` naming the section and value); other
// conditions list what would have fired recently (restocks and
// activity: last 24 h; sale bursts: the current window; catalog:
// announced in the last 7 days). `window` says which. The draft is
// validated against the plan like a new rule (400, 402). Nothing is
// saved or sent.
//
// **API key scope:** `alerts:write`.
func (c *Client) PreviewAlertRule(ctx context.Context, body *AlertRuleInput, opts ...RequestOption) (*AlertPreview, error) {
	var out AlertPreview
	if err := c.do(ctx, http.MethodPost, "/v1/alerts/rules/preview", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAlertRule calls GET /v1/alerts/rules/{ruleId}: One rule.
//
// **API key scope:** `alerts:write`.
func (c *Client) GetAlertRule(ctx context.Context, ruleID int64, opts ...RequestOption) (*AlertRule, error) {
	var out AlertRule
	if err := c.do(ctx, http.MethodGet, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10)), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAlertRule calls PATCH /v1/alerts/rules/{ruleId}: Edit, enable or disable a rule.
//
// Fields given replace the rule's. Changing the scope or condition
// resets its per-event state (the rule re-evaluates from scratch);
// enabling a metric rule evaluates it right away.
//
// **API key scope:** `alerts:write`.
func (c *Client) UpdateAlertRule(ctx context.Context, ruleID int64, body *AlertRulePatch, opts ...RequestOption) (*AlertRule, error) {
	var out AlertRule
	if err := c.do(ctx, http.MethodPatch, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10)), nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAlertRule calls DELETE /v1/alerts/rules/{ruleId}: Delete a rule.
//
// Its firings stay in the history.
//
// **API key scope:** `alerts:write`.
func (c *Client) DeleteAlertRule(ctx context.Context, ruleID int64, opts ...RequestOption) error {
	return c.do(ctx, http.MethodDelete, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10)), nil, nil, nil, opts)
}

// SnoozeAlertRule calls PUT /v1/alerts/rules/{ruleId}/snooze: Snooze a rule.
//
// The rule stays enabled but sends nothing until `until` (at most 90
// days ahead); replaces an earlier snooze. Its matches meanwhile are
// recorded as the skip `snoozed`. A metric or section condition still
// holding when the snooze ends fires then. The snooze simply expires.
//
// **API key scope:** `alerts:write`.
func (c *Client) SnoozeAlertRule(ctx context.Context, ruleID int64, body *AlertRuleSnoozeInput, opts ...RequestOption) (*AlertRule, error) {
	var out AlertRule
	if err := c.do(ctx, http.MethodPut, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10))+"/snooze", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnsnoozeAlertRule calls DELETE /v1/alerts/rules/{ruleId}/snooze: Unsnooze a rule.
//
// Ends the rule's snooze now (a no-op when it has none); an enabled metric or section rule is evaluated right away.
//
// **API key scope:** `alerts:write`.
func (c *Client) UnsnoozeAlertRule(ctx context.Context, ruleID int64, opts ...RequestOption) (*AlertRule, error) {
	var out AlertRule
	if err := c.do(ctx, http.MethodDelete, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10))+"/snooze", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// TestAlertRule calls POST /v1/alerts/rules/{ruleId}/test: Test-fire a rule.
//
// Writes a test firing (marked `test`) about the scope's best current
// match, or any upcoming event in scope, and delivers it to the rule's
// channels and the `alerts` stream topic (not to `alert.fired`
// webhooks), ignoring cooldown and state: use it to check a rule's
// routing end to end. The rule must still be allowed by the org's
// current plan (its scope and features, and its place within the rule
// count; 402); activity it lists stays inside the plan's sales window
// and ended-event bound. An org gets 30 test firings per rolling hour
// (429 with `Retry-After`), and test firings are kept 7 days (real
// firings: the history retention).
//
// **API key scope:** `alerts:write`.
func (c *Client) TestAlertRule(ctx context.Context, ruleID int64, opts ...RequestOption) (*AlertEvent, error) {
	var out AlertEvent
	if err := c.do(ctx, http.MethodPost, "/v1/alerts/rules/"+url.PathEscape(strconv.FormatInt(ruleID, 10))+"/test", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAlertEventsParams holds the query parameters of ListAlertEvents.
type ListAlertEventsParams struct {
	RuleID  *int64
	EventID *int64
	// Firings older than this id (`nextBefore`).
	Before *int64
	Limit  *int
}

func (p *ListAlertEventsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.RuleID != nil {
		q.Set("ruleId", strconv.FormatInt(*p.RuleID, 10))
	}
	if p.EventID != nil {
		q.Set("eventId", strconv.FormatInt(*p.EventID, 10))
	}
	if p.Before != nil {
		q.Set("before", strconv.FormatInt(*p.Before, 10))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListAlertEvents calls GET /v1/alerts/events: Alert history.
//
// The org's firings, newest first, with each channel delivery's status. Snapshot of the `alerts` stream topic.
//
// **API key scope:** `alerts:write`.
func (c *Client) ListAlertEvents(ctx context.Context, params *ListAlertEventsParams, opts ...RequestOption) (*AlertEventPage, error) {
	var out AlertEventPage
	if err := c.do(ctx, http.MethodGet, "/v1/alerts/events", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAlertEventsPages iterates the pages of ListAlertEvents, following nextBefore as before until the last page.
func (c *Client) ListAlertEventsPages(ctx context.Context, params *ListAlertEventsParams, opts ...RequestOption) iter.Seq2[*AlertEventPage, error] {
	var p ListAlertEventsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*AlertEventPage, error) { return c.ListAlertEvents(ctx, &p, opts...) }, func(page *AlertEventPage) bool {
		if page.NextBefore == nil {
			return false
		}
		next := int64(*page.NextBefore)
		p.Before = new(next)
		return true
	})
}

// ListAlertEventsAll iterates every item of ListAlertEvents across pages.
func (c *Client) ListAlertEventsAll(ctx context.Context, params *ListAlertEventsParams, opts ...RequestOption) iter.Seq2[AlertEvent, error] {
	return items(c.ListAlertEventsPages(ctx, params, opts...), func(page *AlertEventPage) []AlertEvent { return page.Items })
}

// ListChannels calls GET /v1/channels: The org's notification channels.
//
// Where alert rules deliver: every channel of the org, oldest first,
// with whether it is verified and enabled and the outcome of its last
// delivery (`lastSentAt`, `lastError`). Route a rule to channels by
// their `id` (`channels` of `POST /v1/alerts/rules`). Every member
// sees them; owners and admins (and their keys) manage them.
//
// **API key scope:** `alerts:write`.
func (c *Client) ListChannels(ctx context.Context, opts ...RequestOption) (*NotificationChannelList, error) {
	var out NotificationChannelList
	if err := c.do(ctx, http.MethodGet, "/v1/channels", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateChannel calls POST /v1/channels: Add a channel.
//
// `discord_webhook`, `slack_webhook`, `teams_webhook`,
// `google_chat_webhook`: the platform's incoming-webhook URL — https,
// on the platform's own host (discord.com/discordapp.com,
// hooks.slack.com, a Teams Workflows URL on *.logic.azure.com or
// *.environment.api.powerplatform.com, chat.googleapis.com) and
// public addresses only; a "connected" message is delivered first and
// a failure rejects the channel (422). `email`: verified at once when
// it is the signed-in caller's own address, else (always for API
// keys) a 6-digit code is emailed to it (`POST …/verify`). `telegram`:
// returns `linkUrl`, a bot deep link that connects the chat it is
// opened in (valid 1 h). Web Push channels are per browser and come
// from the app only.
//
// Limits: an org has at most 25 channels (Web Push not counted; 409).
// Verification messages (email codes, webhook "connected" messages,
// reissued Telegram links) are limited to 10 per hour per org and 5
// per hour per channel and per email address; over a limit the
// answer is 429 with `Retry-After`. Owners and admins only (a key acts
// with its creator's role, capped at admin).
//
// **API key scope:** `alerts:write`.
func (c *Client) CreateChannel(ctx context.Context, body *CreateChannelRequest, opts ...RequestOption) (*NotificationChannelCreated, error) {
	var out NotificationChannelCreated
	if err := c.do(ctx, http.MethodPost, "/v1/channels", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateChannel calls PATCH /v1/channels/{channelId}: Rename, enable or disable a channel.
//
// A disabled channel keeps its rules' routing but receives nothing until enabled again. Owners and admins only.
//
// **API key scope:** `alerts:write`.
func (c *Client) UpdateChannel(ctx context.Context, channelID int64, body *UpdateChannelRequest, opts ...RequestOption) (*NotificationChannel, error) {
	var out NotificationChannel
	if err := c.do(ctx, http.MethodPatch, "/v1/channels/"+url.PathEscape(strconv.FormatInt(channelID, 10)), nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteChannel calls DELETE /v1/channels/{channelId}: Delete a channel.
//
// It is removed from every rule; past deliveries keep its name. Owners and admins only.
//
// **API key scope:** `alerts:write`.
func (c *Client) DeleteChannel(ctx context.Context, channelID int64, opts ...RequestOption) error {
	return c.do(ctx, http.MethodDelete, "/v1/channels/"+url.PathEscape(strconv.FormatInt(channelID, 10)), nil, nil, nil, opts)
}

// TestChannel calls POST /v1/channels/{channelId}/test: Send a test notification now.
//
// Delivered synchronously, within the channel's rate limit; the
// outcome is recorded on the channel (`lastSentAt` or `lastError`) and
// returned. Owners and admins only.
//
// **API key scope:** `alerts:write`.
func (c *Client) TestChannel(ctx context.Context, channelID int64, opts ...RequestOption) (*NotificationChannel, error) {
	var out NotificationChannel
	if err := c.do(ctx, http.MethodPost, "/v1/channels/"+url.PathEscape(strconv.FormatInt(channelID, 10))+"/test", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyChannel calls POST /v1/channels/{channelId}/verify: Verify an email channel with its code.
//
// The 6-digit code emailed to the channel's address. Codes expire after
// 1 h or 5 wrong tries (`POST …/verification` sends a new one). After
// 15 wrong codes over the channel's lifetime (new codes do not reset
// this) the channel cannot be verified any more (409): delete it and
// add the address again. Owners and admins only.
//
// **API key scope:** `alerts:write`.
func (c *Client) VerifyChannel(ctx context.Context, channelID int64, body *VerifyChannelRequest, opts ...RequestOption) (*NotificationChannel, error) {
	var out NotificationChannel
	if err := c.do(ctx, http.MethodPost, "/v1/channels/"+url.PathEscape(strconv.FormatInt(channelID, 10))+"/verify", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// RestartChannelVerification calls POST /v1/channels/{channelId}/verification: Resend an email code or issue a new Telegram link.
//
// Counts against the verification-message limits of `createChannel`
// (per channel, per email address, per org; 429 with `Retry-After`).
// Owners and admins only.
//
// **API key scope:** `alerts:write`.
func (c *Client) RestartChannelVerification(ctx context.Context, channelID int64, opts ...RequestOption) (*NotificationChannelCreated, error) {
	var out NotificationChannelCreated
	if err := c.do(ctx, http.MethodPost, "/v1/channels/"+url.PathEscape(strconv.FormatInt(channelID, 10))+"/verification", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUsageParams holds the query parameters of GetUsage.
type GetUsageParams struct {
	// Start (inclusive); default `to` − 24 h.
	From *time.Time
	// End (exclusive); default now.
	To          *time.Time
	Granularity *GetUsageGranularity
	// Only this key.
	APIKeyID *string
}

func (p *GetUsageParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.From != nil {
		q.Set("from", (*p.From).Format(time.RFC3339Nano))
	}
	if p.To != nil {
		q.Set("to", (*p.To).Format(time.RFC3339Nano))
	}
	if p.Granularity != nil {
		q.Set("granularity", string(*p.Granularity))
	}
	if p.APIKeyID != nil {
		q.Set("apiKeyId", *p.APIKeyID)
	}
	return q
}

// GetUsage calls GET /v1/usage: API usage and quotas of the org.
//
// Requests (by endpoint group), 429s, server errors and stream
// connections/messages of the org's API keys, metered per key and
// minute, bucketed by `granularity` over [`from`, `to`) (default: the
// last 24 hours by hour; minute buckets cover at most 6 hours, hour
// buckets 7 days, day buckets 90 days). The current minute is included
// live. `quota` is the plan's limits; `current` the live usage against
// them. Any key may read it (no scope needed).
func (c *Client) GetUsage(ctx context.Context, params *GetUsageParams, opts ...RequestOption) (*Usage, error) {
	var out Usage
	if err := c.do(ctx, http.MethodGet, "/v1/usage", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListWebhooks calls GET /v1/webhooks: The org's webhook endpoints.
//
// **API key scope:** `webhooks:write`.
func (c *Client) ListWebhooks(ctx context.Context, opts ...RequestOption) (*WebhookList, error) {
	var out WebhookList
	if err := c.do(ctx, http.MethodGet, "/v1/webhooks", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateWebhook calls POST /v1/webhooks: Register a webhook endpoint.
//
// Deliveries are signed POSTs (see the `Webhooks` guide): header
// `Status302-Signature: t=<unix>,v1=<hex HMAC-SHA256(secret, t + "." + body)>`,
// plus `Status302-Delivery-Id` (idempotency key) and `Status302-Event`. The
// signing secret is returned once here (and on rotate). The URL must be
// https and resolve to public addresses only. `list.sale`,
// `event.status_changed` and the market-activity types
// (`inventory.restocked`, `inventory.sold`, `resale.listed`,
// `resale.sold`, `marketplace.listed`, `marketplace.sold`) need a
// `listId`: they cover the events of that list. `markets` narrows
// `list.sale` and the market-activity types to some markets. Needs a
// plan with webhooks (Elite). Audited.
//
// **API key scope:** `webhooks:write`.
func (c *Client) CreateWebhook(ctx context.Context, body *WebhookInput, opts ...RequestOption) (*WebhookSecret, error) {
	var out WebhookSecret
	if err := c.do(ctx, http.MethodPost, "/v1/webhooks", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWebhook calls GET /v1/webhooks/{webhookId}: One webhook endpoint.
//
// **API key scope:** `webhooks:write`.
func (c *Client) GetWebhook(ctx context.Context, webhookID string, opts ...RequestOption) (*Webhook, error) {
	var out Webhook
	if err := c.do(ctx, http.MethodGet, "/v1/webhooks/"+url.PathEscape(webhookID), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateWebhook calls PATCH /v1/webhooks/{webhookId}: Change an endpoint, or re-enable a disabled one.
//
// Fields given replace the endpoint's; `status=active` re-enables it and resets its failure count. Audited.
//
// **API key scope:** `webhooks:write`.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, body *WebhookPatch, opts ...RequestOption) (*Webhook, error) {
	var out Webhook
	if err := c.do(ctx, http.MethodPatch, "/v1/webhooks/"+url.PathEscape(webhookID), nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteWebhook calls DELETE /v1/webhooks/{webhookId}: Delete an endpoint.
//
// Pending deliveries are dropped. Audited.
//
// **API key scope:** `webhooks:write`.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string, opts ...RequestOption) error {
	return c.do(ctx, http.MethodDelete, "/v1/webhooks/"+url.PathEscape(webhookID), nil, nil, nil, opts)
}

// RotateWebhookSecret calls POST /v1/webhooks/{webhookId}/rotate-secret: Replace an endpoint's signing secret.
//
// The new secret signs every delivery from now on (including retries) and is returned once. Audited.
//
// **API key scope:** `webhooks:write`.
func (c *Client) RotateWebhookSecret(ctx context.Context, webhookID string, opts ...RequestOption) (*WebhookSecret, error) {
	var out WebhookSecret
	if err := c.do(ctx, http.MethodPost, "/v1/webhooks/"+url.PathEscape(webhookID)+"/rotate-secret", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// TestWebhook calls POST /v1/webhooks/{webhookId}/test: Send a test event.
//
// Queues one delivery of a sample `type` event (default: the first
// subscribed type) marked `"test": true` in the envelope. Disabled
// endpoints are tested too (to check a fix before re-enabling).
//
// **API key scope:** `webhooks:write`.
func (c *Client) TestWebhook(ctx context.Context, webhookID string, body *WebhookTestInput, opts ...RequestOption) (*WebhookDelivery, error) {
	var out WebhookDelivery
	if err := c.do(ctx, http.MethodPost, "/v1/webhooks/"+url.PathEscape(webhookID)+"/test", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListWebhookDeliveriesParams holds the query parameters of ListWebhookDeliveries.
type ListWebhookDeliveriesParams struct {
	Status *WebhookDeliveryStatus
	Cursor *string
	Limit  *int
}

func (p *ListWebhookDeliveriesParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Status != nil {
		q.Set("status", string(*p.Status))
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListWebhookDeliveries calls GET /v1/webhooks/{webhookId}/deliveries: Delivery log of an endpoint.
//
// Newest first, cursor-paginated; kept 30 days.
//
// **API key scope:** `webhooks:write`.
func (c *Client) ListWebhookDeliveries(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) (*WebhookDeliveryPage, error) {
	var out WebhookDeliveryPage
	if err := c.do(ctx, http.MethodGet, "/v1/webhooks/"+url.PathEscape(webhookID)+"/deliveries", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListWebhookDeliveriesPages iterates the pages of ListWebhookDeliveries, following nextCursor as cursor until the last page.
func (c *Client) ListWebhookDeliveriesPages(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) iter.Seq2[*WebhookDeliveryPage, error] {
	var p ListWebhookDeliveriesParams
	if params != nil {
		p = *params
	}
	return pages(func() (*WebhookDeliveryPage, error) { return c.ListWebhookDeliveries(ctx, webhookID, &p, opts...) }, func(page *WebhookDeliveryPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListWebhookDeliveriesAll iterates every item of ListWebhookDeliveries across pages.
func (c *Client) ListWebhookDeliveriesAll(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) iter.Seq2[WebhookDelivery, error] {
	return items(c.ListWebhookDeliveriesPages(ctx, webhookID, params, opts...), func(page *WebhookDeliveryPage) []WebhookDelivery { return page.Items })
}

// RedeliverWebhook calls POST /v1/webhooks/{webhookId}/deliveries/{deliveryId}/redeliver: Deliver a past event again.
//
// Queues a new delivery of the same envelope (same `id`, a new delivery log row with `redeliveryOf`).
//
// **API key scope:** `webhooks:write`.
func (c *Client) RedeliverWebhook(ctx context.Context, webhookID string, deliveryID string, opts ...RequestOption) (*WebhookDelivery, error) {
	var out WebhookDelivery
	if err := c.do(ctx, http.MethodPost, "/v1/webhooks/"+url.PathEscape(webhookID)+"/deliveries/"+url.PathEscape(deliveryID)+"/redeliver", nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListExportsParams holds the query parameters of ListExports.
type ListExportsParams struct {
	Cursor *string
	Limit  *int
}

func (p *ListExportsParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Cursor != nil {
		q.Set("cursor", *p.Cursor)
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return q
}

// ListExports calls GET /v1/exports: The org's exports.
//
// Newest first, cursor-paginated; files are deleted 7 days after they were written.
//
// **API key scope:** `catalog:read`.
func (c *Client) ListExports(ctx context.Context, params *ListExportsParams, opts ...RequestOption) (*ExportPage, error) {
	var out ExportPage
	if err := c.do(ctx, http.MethodGet, "/v1/exports", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListExportsPages iterates the pages of ListExports, following nextCursor as cursor until the last page.
func (c *Client) ListExportsPages(ctx context.Context, params *ListExportsParams, opts ...RequestOption) iter.Seq2[*ExportPage, error] {
	var p ListExportsParams
	if params != nil {
		p = *params
	}
	return pages(func() (*ExportPage, error) { return c.ListExports(ctx, &p, opts...) }, func(page *ExportPage) bool {
		if page.NextCursor == nil {
			return false
		}
		next := string(*page.NextCursor)
		if next == "" {
			return false
		}
		p.Cursor = new(next)
		return true
	})
}

// ListExportsAll iterates every item of ListExports across pages.
func (c *Client) ListExportsAll(ctx context.Context, params *ListExportsParams, opts ...RequestOption) iter.Seq2[Export, error] {
	return items(c.ListExportsPages(ctx, params, opts...), func(page *ExportPage) []Export { return page.Items })
}

// CreateExport calls POST /v1/exports: Start an async export.
//
// Writes every row of a REST read to a CSV or JSON Lines file in the
// background; poll `GET /v1/exports/{exportId}` until `status` is
// `succeeded`, then download `downloadUrl` (signed, valid 15 minutes,
// no credentials needed). `filters` are the query parameters of the read
// the kind mirrors, with the same validation and plan clamps:
//
// | kind | read | required filters |
// |---|---|---|
// | `events` | `GET /v1/events` | – |
// | `sales` | `GET /v1/sales/feed`, or `GET /v1/events/{eventId}/sales` with `eventId` | – |
// | `series` | `GET /v1/events/{eventId}/series` | `eventId` |
// | `list` | `GET /v1/lists/{listId}/grid` (the list's rows with their metrics) | `listId` |
//
// Scopes: `catalog:read` for events and lists, `sales:read` for sales,
// `market:read` for series. Rows are capped by the plan's export limit
// and, on every plan, by 1,000,000 rows and 256 MiB per file
// (`truncated` tells). An org may have 3 exports queued or running;
// another is answered 429 (`limit` `concurrent`) with Retry-After.
//
// **API key scope:** `catalog:read`.
func (c *Client) CreateExport(ctx context.Context, body *ExportInput, opts ...RequestOption) (*Export, error) {
	var out Export
	if err := c.do(ctx, http.MethodPost, "/v1/exports", nil, jsonBody(body), &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetExport calls GET /v1/exports/{exportId}: Status of an export.
//
// **API key scope:** `catalog:read`.
func (c *Client) GetExport(ctx context.Context, exportID string, opts ...RequestOption) (*Export, error) {
	var out Export
	if err := c.do(ctx, http.MethodGet, "/v1/exports/"+url.PathEscape(exportID), nil, nil, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadExportParams holds the query parameters of DownloadExport.
type DownloadExportParams struct {
	// Required.
	Expires int64
	// Required.
	Signature string
}

func (p *DownloadExportParams) values() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	q.Set("expires", strconv.FormatInt(p.Expires, 10))
	q.Set("signature", p.Signature)
	return q
}

// DownloadExport calls GET /v1/exports/{exportId}/download: Download an export file (signed URL).
//
// The `downloadUrl` of a succeeded export when the object store cannot
// presign URLs itself: authenticated by its `expires` and `signature`
// query parameters, not by credentials.
func (c *Client) DownloadExport(ctx context.Context, exportID string, params *DownloadExportParams, opts ...RequestOption) ([]byte, error) {
	var out []byte
	if err := c.do(ctx, http.MethodGet, "/v1/exports/"+url.PathEscape(exportID)+"/download", params.values(), nil, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}
