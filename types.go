// Code generated from the Status302 public OpenAPI spec (https://status302.com/docs/api/openapi.yaml). DO NOT EDIT.

package status302

import "time"

// ListEventsSort is a schema of the public API.
type ListEventsSort string

// Values of ListEventsSort.
const (
	ListEventsSortStartsAt     ListEventsSort = "startsAt"
	ListEventsSortFirstSeen    ListEventsSort = "firstSeen"
	ListEventsSortOnsale       ListEventsSort = "onsale"
	ListEventsSortFloor        ListEventsSort = "floor"
	ListEventsSortChg24h       ListEventsSort = "chg24h"
	ListEventsSortChg72h       ListEventsSort = "chg72h"
	ListEventsSortSales24h     ListEventsSort = "sales24h"
	ListEventsSortDepletion24h ListEventsSort = "depletion24h"
	ListEventsSortMovement     ListEventsSort = "movement"
	ListEventsSortInventory    ListEventsSort = "inventory"
	ListEventsSortSellThrough  ListEventsSort = "sellThrough"
	ListEventsSortSpread       ListEventsSort = "spread"
	ListEventsSortFloorPair    ListEventsSort = "floorPair"
	ListEventsSortMedian       ListEventsSort = "median"
	ListEventsSortAvg24h       ListEventsSort = "avg24h"
	ListEventsSortAvg72h       ListEventsSort = "avg72h"
	ListEventsSortAvg7d        ListEventsSort = "avg7d"
	ListEventsSortAvg30d       ListEventsSort = "avg30d"
	ListEventsSortChg7d        ListEventsSort = "chg7d"
	ListEventsSortSales72h     ListEventsSort = "sales72h"
	ListEventsSortVolume24h    ListEventsSort = "volume24h"
	ListEventsSortSaleAvg24h   ListEventsSort = "saleAvg24h"
	ListEventsSortSaleAvg72h   ListEventsSort = "saleAvg72h"
	ListEventsSortVelocity     ListEventsSort = "velocity"
	ListEventsSortChurn24h     ListEventsSort = "churn24h"
	ListEventsSortSelloutHours ListEventsSort = "selloutHours"
)

// ListEventsOrder is a schema of the public API.
type ListEventsOrder string

// Values of ListEventsOrder.
const (
	ListEventsOrderAsc  ListEventsOrder = "asc"
	ListEventsOrderDesc ListEventsOrder = "desc"
)

// ListEventListingsMarket is a schema of the public API.
type ListEventListingsMarket string

// Values of ListEventListingsMarket.
const (
	ListEventListingsMarketPrimary   ListEventListingsMarket = "primary"
	ListEventListingsMarketTmResale  ListEventListingsMarket = "tm_resale"
	ListEventListingsMarketSecondary ListEventListingsMarket = "secondary"
)

// ListEventListingsSort is a schema of the public API.
type ListEventListingsSort string

// Values of ListEventListingsSort.
const (
	ListEventListingsSortPrice     ListEventListingsSort = "price"
	ListEventListingsSortPriceDesc ListEventListingsSort = "-price"
)

// ListMoversWindow is a schema of the public API.
type ListMoversWindow string

// Values of ListMoversWindow.
const (
	ListMoversWindow24h ListMoversWindow = "24h"
	ListMoversWindow72h ListMoversWindow = "72h"
	ListMoversWindow7d  ListMoversWindow = "7d"
)

// ListMoversDirection is a schema of the public API.
type ListMoversDirection string

// Values of ListMoversDirection.
const (
	ListMoversDirectionUp   ListMoversDirection = "up"
	ListMoversDirectionDown ListMoversDirection = "down"
)

// GetCompareSeriesAlign is a schema of the public API.
type GetCompareSeriesAlign string

// Values of GetCompareSeriesAlign.
const (
	GetCompareSeriesAlignCalendar    GetCompareSeriesAlign = "calendar"
	GetCompareSeriesAlignDaysToEvent GetCompareSeriesAlign = "days_to_event"
)

// ListCalendarItemsSort is a schema of the public API.
type ListCalendarItemsSort string

// Values of ListCalendarItemsSort.
const (
	ListCalendarItemsSortTime       ListCalendarItemsSort = "time"
	ListCalendarItemsSortProminence ListCalendarItemsSort = "prominence"
)

// GetUsageGranularity is a schema of the public API.
type GetUsageGranularity string

// Values of GetUsageGranularity.
const (
	GetUsageGranularityMinute GetUsageGranularity = "minute"
	GetUsageGranularityHour   GetUsageGranularity = "hour"
	GetUsageGranularityDay    GetUsageGranularity = "day"
)

// StreamTicket is a schema of the public API.
type StreamTicket struct {
	// Opaque single-use credential; pass as `?ticket=` on `GET /v1/stream`.
	Ticket    string    `json:"ticket"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// CalendarKind: `onsale`: public on-sale start; `presale`: presale start; `event`: show start.
type CalendarKind string

// Values of CalendarKind.
const (
	CalendarKindOnsale  CalendarKind = "onsale"
	CalendarKindPresale CalendarKind = "presale"
	CalendarKindEvent   CalendarKind = "event"
)

// CalendarCounts is a schema of the public API.
type CalendarCounts struct {
	Onsale  int `json:"onsale"`
	Presale int `json:"presale"`
	Event   int `json:"event"`
}

// CalendarDays is a schema of the public API.
type CalendarDays struct {
	TZ   string        `json:"tz"`
	From string        `json:"from"`
	To   string        `json:"to"`
	Days []CalendarDay `json:"days"`
}

// CalendarDay is a schema of the public API.
type CalendarDay struct {
	Date   string         `json:"date"`
	Counts CalendarCounts `json:"counts"`
	// Highest on-sale heat (`CalendarItem.heat`) among the day's on-sale and
	// presale items matching the filters; null when none is scored. Drives
	// the calendar's per-cell heat colouring.
	Heat *float64 `json:"heat"`
	// The day's most prominent items (ties by time; one item per event first).
	Top []CalendarItem `json:"top"`
}

// CalendarItemPage is a schema of the public API.
type CalendarItemPage struct {
	Items []CalendarItem `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
}

// CalendarItem is a schema of the public API.
type CalendarItem struct {
	// Stable item key (`kind:eventId:unixSeconds`).
	ID   string       `json:"id"`
	Kind CalendarKind `json:"kind"`
	At   time.Time    `json:"at"`
	// End of the sale window (the latest end of a presale stack); null for shows and open-ended windows.
	EndsAt *time.Time `json:"endsAt"`
	// Local day of `at` in the requested `tz`.
	Date string `json:"date"`
	// The presales starting together, by name (`presale` items only; empty otherwise).
	Presales []CalendarPresale `json:"presales"`
	// 0–100 ranking prior from catalog facts: 25 if hosted on
	// Ticketmaster, up to 40 for venue capacity (ln(1+capacity)/ln(1+80000),
	// see `CalendarEvent.capacity`), up to 20 for the event's presale count
	// (4+ = full), up to 15 for the advertised face-value maximum
	// (ln(1+max)/ln(1+1000), major units, not currency-converted).
	// Cancelled events score 0.
	Prominence float64 `json:"prominence"`
	// The event's on-sale heat, 0–100, recomputed on every metrics refresh
	// for listed upcoming events with an on-sale or presale
	// starting in the last 7 days or later; null for other events. A weighted
	// mean of component scores (weights renormalised over the known ones):
	// 30% the median secondary spread of the reference events (+100% = full),
	// 20% their median primary sell-through, 15% their median secondary
	// median ÷ the advertised face maximum (0.5× = 0, 2× = full), 10% venue
	// capacity (log scale, 80k = full), 10% routing (1 / (1 + log2 shows of
	// the headliner within ±90 days)), 15% presales (4+ = full). Reference
	// events are the headliner's other events with metrics (2+), else the
	// segment's (5+), else every event's.
	Heat  *float64      `json:"heat"`
	Event CalendarEvent `json:"event"`
}

// CalendarPresale is a schema of the public API.
type CalendarPresale struct {
	Name   string     `json:"name"`
	EndsAt *time.Time `json:"endsAt"`
}

// CalendarEvent is a schema of the public API.
type CalendarEvent struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	TmEventID *string `json:"tmEventId"`
	// The event's primary seller (`EventSummary.primarySource`).
	PrimarySource *PrimarySource `json:"primarySource"`
	// Ticketmaster (or other primary) event page.
	URL       *string     `json:"url"`
	StartsAt  *time.Time  `json:"startsAt"`
	LocalDate *string     `json:"localDate"`
	Timezone  *string     `json:"timezone"`
	Status    EventStatus `json:"status"`
	Venue     *VenueRef   `json:"venue,omitempty"`
	// The event's attractions in billing order.
	Attractions []AttractionRef `json:"attractions"`
	// A compiled seat map is linked to the event.
	HasSeatMap bool `json:"hasSeatMap"`
	// Declared places of the compiled seat map: seats plus GA capacity (the larger of
	// the seat count and the sections' declared totals). A venue capacity proxy; null
	// without a map or when the map declares none.
	Capacity *int `json:"capacity"`
}

// EventMapStatus is a schema of the public API.
type EventMapStatus string

// Values of EventMapStatus.
const (
	EventMapStatusOk      EventMapStatus = "ok"
	EventMapStatusNoMap   EventMapStatus = "no_map"
	EventMapStatusError   EventMapStatus = "error"
	EventMapStatusPending EventMapStatus = "pending"
)

// EventMap is a schema of the public API.
type EventMap struct {
	EventID int64          `json:"eventId"`
	Status  EventMapStatus `json:"status"`
	// Ticketmaster venue configuration of the event's map.
	VenueConfigID *string `json:"venueConfigId"`
	// Last map fetch; null while pending.
	FetchedAt *time.Time `json:"fetchedAt"`
	// Why the last fetched map could not be compiled (status `error`).
	Error *string `json:"error"`
	// The inventory manifest version current when the map was fetched.
	ManifestVersion *string  `json:"manifestVersion"`
	Seatmap         *Seatmap `json:"seatmap,omitempty"`
}

// Seatmap: A compiled artifact (absent for `pending`, `no_map`, and errors without an earlier map).
type Seatmap struct {
	// Lowercase hex SHA-256 of the uncompressed artifact.
	Hash string `json:"hash"`
	// Path of the artifact (`/v1/seatmaps/{hash}`).
	URL string `json:"url"`
	// Absolute URL of the same artifact on the public map artifact
	// store (a CDN in front of the object storage bucket,
	// `…/maps/{venueConfigId}/{hash}.tmap`), present when the
	// deployment publishes artifacts there. Always brotli-encoded
	// (`Content-Encoding: br`), immutable, readable cross-origin
	// without credentials. Prefer it to `url` when present.
	CdnURL *string `json:"cdnUrl,omitempty"`
	// .tmap format version (2; format 1 artifacts were deleted and are no longer served).
	Format        int    `json:"format"`
	VenueConfigID string `json:"venueConfigId"`
	PageWidth     int    `json:"pageWidth"`
	PageHeight    int    `json:"pageHeight"`
	Sections      int    `json:"sections"`
	GaSections    int    `json:"gaSections"`
	Places        int    `json:"places"`
	RawBytes      int    `json:"rawBytes"`
	// Size with brotli (`Content-Encoding: br`).
	CompressedBytes int `json:"compressedBytes"`
}

// MapOverlay is a schema of the public API.
type MapOverlay struct {
	EventID int64 `json:"eventId"`
	// The artifact the indexes refer to.
	SeatmapHash string `json:"seatmapHash"`
	// Time of the inventory snapshot.
	AsOf time.Time `json:"asOf"`
	// Availability revision of the snapshot.
	Revision int64 `json:"revision"`
	// Places in the artifact (bits per bitset, entries in `priceLevel`).
	Places int `json:"places"`
	// Bitset of places open on the primary market.
	Primary []byte `json:"primary"`
	// Bitset of places listed on TM resale (the tmresale listings that name their places).
	Resale []byte `json:"resale"`
	// Little-endian int16 per place; index into `priceLevels`, -1 when unpriced.
	PriceLevel []byte `json:"priceLevel"`
	// Price levels in use, cheapest all-in first.
	PriceLevels []OverlayPriceLevel `json:"priceLevels"`
	// Sections with open stock, by artifact section index.
	Sections []OverlaySection `json:"sections"`
	// Open primary GA stock per GA section.
	Ga          []OverlayGA        `json:"ga"`
	Currency    *string            `json:"currency"`
	Diagnostics OverlayDiagnostics `json:"diagnostics"`
	// Per-scope section values; present when `mode` is given.
	SectionLayers []OverlaySectionLayer `json:"sectionLayers,omitempty"`
}

// OverlayMode is a schema of the public API.
type OverlayMode string

// Values of OverlayMode.
const (
	OverlayModeAvailability OverlayMode = "availability"
	OverlayModePrice        OverlayMode = "price"
	OverlayModeDelta        OverlayMode = "delta"
	OverlayModeSales        OverlayMode = "sales"
	OverlayModeSpread       OverlayMode = "spread"
)

// OverlayMarket is a schema of the public API.
type OverlayMarket string

// Values of OverlayMarket.
const (
	OverlayMarketPrimary   OverlayMarket = "primary"
	OverlayMarketTmResale  OverlayMarket = "tm_resale"
	OverlayMarketSecondary OverlayMarket = "secondary"
)

// OverlaySectionLayer is a schema of the public API.
type OverlaySectionLayer struct {
	// `primary`, `tm_resale`, `secondary:<source>` or `secondary:best`.
	Scope string `json:"scope"`
	// Currency of the layer's prices.
	Currency *string `json:"currency"`
	// Latest change of the layer's sections.
	AsOf *time.Time `json:"asOf"`
	// Sections on the map with stock or (mode `sales`, `spread`) recent sales, by index.
	Sections []OverlaySectionValue `json:"sections"`
}

// OverlaySectionValue: One section's values in a layer; the mode decides which optional fields
// are present. With map prices every mode but `availability` carries
// `floorMinor` and `floorChg24h` (the map's hover shows them in every
// mode) besides its own: `price` adds `medianMinor`, `sales` adds
// `sales24h`, `sales72h` and `saleAvg24hMinor`, `spread` adds `spread`
// and `roi`. Without map prices only `availability` and `sales` (counts
// only) are available.
type OverlaySectionValue struct {
	// Artifact section index.
	Index       int             `json:"index"`
	Key         string          `json:"key"`
	Tickets     int             `json:"tickets"`
	FloorMinor  Nullable[int64] `json:"floorMinor,omitzero"`
	MedianMinor Nullable[int64] `json:"medianMinor,omitzero"`
	// (floor − floor 24h ago) ÷ floor 24h ago.
	FloorChg24h     Nullable[float64] `json:"floorChg24h,omitzero"`
	Sales24h        *int              `json:"sales24h,omitempty"`
	Sales72h        *int              `json:"sales72h,omitempty"`
	SaleAvg24hMinor Nullable[int64]   `json:"saleAvg24hMinor,omitzero"`
	Spread          Nullable[float64] `json:"spread,omitzero"`
	ROI             Nullable[float64] `json:"roi,omitzero"`
	// Marketplace layers only: the section's tickets by the row their
	// listings name, by row label; their tickets add up to `tickets`.
	// Listings naming no row sit in the entry with an empty `row`.
	Rows []OverlaySeatRow `json:"rows,omitempty"`
}

// OverlaySeatRow: A marketplace's tickets in one row of a section. `row` is the row label
// as the listings give it (match it against the artifact's row names
// folded to letters and digits, lower-cased), empty for listings naming
// no row. `seated` of the tickets come from listings naming their seats:
// by seat label range (`seats`) or by artifact place index (`places`).
type OverlaySeatRow struct {
	Row     string `json:"row"`
	Tickets int    `json:"tickets"`
	// Tickets of listings with seat positions (`seats` or `places`).
	Seated int `json:"seated"`
	// Seat label ranges of the listings naming their seats (inclusive, as listed), by label.
	Seats []OverlaySeatRange `json:"seats"`
	// Artifact place indexes of the listings that came with them, ascending.
	Places []int `json:"places"`
}

// OverlaySeatRange is a schema of the public API.
type OverlaySeatRange struct {
	// First seat label.
	From string `json:"from"`
	// Last seat label (equal to `from` for one seat).
	To string `json:"to"`
}

// EventSections is a schema of the public API.
type EventSections struct {
	EventID int64  `json:"eventId"`
	Scope   string `json:"scope"`
	// Currency of every price of the table.
	Currency *string `json:"currency"`
	// Latest change of the scope's sections.
	AsOf *time.Time `json:"asOf"`
	// The seat map `index` refers to.
	SeatmapHash *string `json:"seatmapHash"`
	// Sections in map order (sections without an index last, by key).
	Items []SectionRow `json:"items"`
}

// SectionRow is a schema of the public API.
type SectionRow struct {
	SectionKey string `json:"sectionKey"`
	// Artifact section index; null when the key names no section of the map.
	Index *int `json:"index"`
	// Tickets listed now (0 for a section that only sold recently).
	Tickets     int    `json:"tickets"`
	FloorMinor  *int64 `json:"floorMinor"`
	MedianMinor *int64 `json:"medianMinor"`
	// When the section's values last changed.
	ChangedAt *time.Time `json:"changedAt"`
	// (floor − floor 24h ago) ÷ floor 24h ago; null without both.
	FloorChg24h *float64 `json:"floorChg24h"`
	// Tickets sold in the last 24 hours.
	Sales24h int `json:"sales24h"`
	Sales72h int `json:"sales72h"`
	// Quantity-weighted mean sale price over 24 hours.
	SaleAvg24hMinor *int64 `json:"saleAvg24hMinor"`
	// Floor (`primary`, the best marketplace floor) ÷ the section's primary floor − 1.
	Spread *float64 `json:"spread"`
	// (median resale price over 72h × (1 − sellerFee) − primary floor) ÷ primary floor.
	ROI *float64 `json:"roi"`
}

// SectionSeries is a schema of the public API.
type SectionSeries struct {
	EventID    int64            `json:"eventId"`
	SectionKey string           `json:"sectionKey"`
	Scope      string           `json:"scope"`
	Res        SeriesResolution `json:"res"`
	From       time.Time        `json:"from"`
	To         time.Time        `json:"to"`
	Currency   *string          `json:"currency"`
	T          []time.Time      `json:"t"`
	Floor      SeriesColumn     `json:"floor"`
	Median     SeriesColumn     `json:"median"`
	Tickets    SeriesColumn     `json:"tickets"`
	Sales      SeriesColumn     `json:"sales"`
}

// OverlayPriceLevel is a schema of the public API.
type OverlayPriceLevel struct {
	// AVSC price key.
	Key       int32  `json:"key"`
	OfferName string `json:"offerName"`
	ListMinor int64  `json:"listMinor"`
	// All-in price per ticket.
	TotalMinor int64  `json:"totalMinor"`
	Currency   string `json:"currency"`
	// Dynamic (seat-level, e.g. Platinum) pricing.
	SeatLevel bool `json:"seatLevel"`
}

// PriceLevelHistory is a schema of the public API.
type PriceLevelHistory struct {
	EventID int64               `json:"eventId"`
	Items   []PriceLevelVersion `json:"items"`
}

// PriceLevelVersion: One version of a Ticketmaster primary price level. Money is minor units of
// `currency` per ticket: `faceMinor` the face value, `listMinor` the list
// price, the fee parts, and `totalMinor` the all-in price.
type PriceLevelVersion struct {
	// AVSC price key (`OverlayPriceLevel.key`).
	PriceKey      int32  `json:"priceKey"`
	PriceID       int32  `json:"priceId"`
	Zone          int32  `json:"zone"`
	OfferID       string `json:"offerId"`
	OfferName     string `json:"offerName"`
	Currency      string `json:"currency"`
	ListMinor     int64  `json:"listMinor"`
	FaceMinor     int64  `json:"faceMinor"`
	FacilityMinor int64  `json:"facilityMinor"`
	ServiceMinor  int64  `json:"serviceMinor"`
	DistanceMinor int64  `json:"distanceMinor"`
	TaxMinor      int64  `json:"taxMinor"`
	TotalMinor    int64  `json:"totalMinor"`
	// Dynamic (seat-level, e.g. Platinum) pricing.
	SeatLevel bool      `json:"seatLevel"`
	ValidFrom time.Time `json:"validFrom"`
	// When the version was replaced or the level disappeared; null while in use.
	ValidTo *time.Time `json:"validTo"`
}

// TicketOfferList is a schema of the public API.
type TicketOfferList struct {
	EventID int64 `json:"eventId"`
	// Currency of `minMinor` and `maxMinor`; null without priced offers or without map prices.
	Currency *string `json:"currency"`
	// Offers with price levels in use first, then by name.
	Items []TicketOffer `json:"items"`
}

// TicketOffer: One Ticketmaster offer (ticket type, e.g. `ADULT`, `PLAT1`) over its
// primary price levels.
type TicketOffer struct {
	// The levels' `offerName`; the series' `offer`.
	Name string `json:"name"`
	// Price levels of the offer in use now (0 for an offer no longer sold).
	Levels int32 `json:"levels"`
	// Open primary seats and GA stock at those levels; null before the event's first poll.
	Tickets *int64 `json:"tickets"`
	// Cheapest all-in price of those levels; null without levels in use or without map prices.
	MinMinor *int64 `json:"minMinor"`
	// Dearest all-in price of those levels; null like `minMinor`.
	MaxMinor *int64 `json:"maxMinor"`
}

// TicketTypeList is a schema of the public API.
type TicketTypeList struct {
	EventID int64 `json:"eventId"`
	// The event's primary seller (`EventSummary.primarySource`).
	Source *PrimarySource `json:"source"`
	// Current ticket types in the source's order, then removed ones (most recently removed first).
	Items []TicketType `json:"items"`
	// Changes, newest first, at most 200.
	Changes []TicketTypeChange `json:"changes"`
}

// TicketType: One ticket type of a primary source's event: general admission at one
// price, without a remaining count. Money is minor units of `currency` per
// ticket: `priceMinor` before tax, `allInMinor` what a buyer pays.
type TicketType struct {
	// The source's ticket type id.
	ExternalID string `json:"externalId"`
	Name       string `json:"name"`
	// `on_sale`, `sold_out`, `off_sale`, or another source status in snake case.
	Status string `json:"status"`
	// The source runs a waiting list for it.
	Waitlist bool `json:"waitlist"`
	// Price tier (release) index; -1 when the source names no tier.
	TierIndex int32 `json:"tierIndex"`
	// Price tier name, e.g. "Second release".
	TierName    *string   `json:"tierName"`
	PriceMinor  int64     `json:"priceMinor"`
	AllInMinor  int64     `json:"allInMinor"`
	Currency    string    `json:"currency"`
	MaxPerOrder *int32    `json:"maxPerOrder"`
	FirstSeen   time.Time `json:"firstSeen"`
	LastChanged time.Time `json:"lastChanged"`
	// When the source stopped listing the type; null while listed.
	RemovedAt *time.Time `json:"removedAt"`
}

// TicketTypeChangeKind: `added` / `removed`: the type appeared in or left the list; `status`: its
// status changed; `price`: its all-in price or tier changed; `waitlist`: its
// waiting list opened or closed.
type TicketTypeChangeKind string

// Values of TicketTypeChangeKind.
const (
	TicketTypeChangeKindAdded    TicketTypeChangeKind = "added"
	TicketTypeChangeKindRemoved  TicketTypeChangeKind = "removed"
	TicketTypeChangeKindStatus   TicketTypeChangeKind = "status"
	TicketTypeChangeKindPrice    TicketTypeChangeKind = "price"
	TicketTypeChangeKindWaitlist TicketTypeChangeKind = "waitlist"
)

// TicketTypeChange: One change of a ticket type; the fields of other kinds are null.
type TicketTypeChange struct {
	Ts             time.Time            `json:"ts"`
	ExternalID     string               `json:"externalId"`
	Kind           TicketTypeChangeKind `json:"kind"`
	Name           string               `json:"name"`
	PrevStatus     *string              `json:"prevStatus"`
	Status         *string              `json:"status"`
	PrevTierName   *string              `json:"prevTierName"`
	TierName       *string              `json:"tierName"`
	PrevAllInMinor *int64               `json:"prevAllInMinor"`
	AllInMinor     *int64               `json:"allInMinor"`
	Currency       *string              `json:"currency"`
	PrevWaitlist   *bool                `json:"prevWaitlist"`
	Waitlist       *bool                `json:"waitlist"`
}

// PrimarySource: An event's primary seller: `tm` (Ticketmaster-hosted) or a primary source with its own catalog (`dice`).
type PrimarySource string

// Values of PrimarySource.
const (
	PrimarySourceTm   PrimarySource = "tm"
	PrimarySourceDice PrimarySource = "dice"
)

// OverlaySection is a schema of the public API.
type OverlaySection struct {
	Index int `json:"index"`
	// Open primary seats plus GA stock.
	PrimaryOpen int `json:"primaryOpen"`
	// Tickets of the section's current TM resale listings.
	ResaleOpen int `json:"resaleOpen"`
	// Cheapest all-in primary price; null when nothing open is priced.
	PrimaryFloorMinor *int64 `json:"primaryFloorMinor"`
	// Cheapest all-in TM resale listing price per ticket; null without one.
	ResaleFloorMinor *int64 `json:"resaleFloorMinor"`
}

// OverlayGA is a schema of the public API.
type OverlayGA struct {
	Section    int              `json:"section"`
	Open       int              `json:"open"`
	FloorMinor *int64           `json:"floorMinor"`
	Levels     []OverlayGALevel `json:"levels"`
}

// OverlayGALevel is a schema of the public API.
type OverlayGALevel struct {
	// Index into `priceLevels`, -1 when unpriced.
	Level int `json:"level"`
	Open  int `json:"open"`
}

// OverlayDiagnostics: How the inventory joined the map.
type OverlayDiagnostics struct {
	// Open primary seats in the snapshot.
	PrimarySeats int `json:"primarySeats"`
	// PlaceIds named by the current TM resale listings.
	ResaleSeats int `json:"resaleSeats"`
	// Of them, placed on the map (bits set in `primary`).
	PrimaryMapped int `json:"primaryMapped"`
	// Of them, placed on the map (bits set in `resale`).
	ResaleMapped int `json:"resaleMapped"`
	// Open seats (either market) whose manifest entry has no placeId.
	SeatsWithoutPlaceID int `json:"seatsWithoutPlaceId"`
	// Open seats whose placeId the map does not contain.
	PlaceIDsNotInMap int `json:"placeIdsNotInMap"`
	GaOpen           int `json:"gaOpen"`
	// GA stock in manifest sections no map section matched.
	GaUnmatched         int      `json:"gaUnmatched"`
	UnmatchedGASections []string `json:"unmatchedGASections"`
	ManifestSeats       int      `json:"manifestSeats"`
	ManifestPlaceIDs    int      `json:"manifestPlaceIds"`
	// Manifest seats whose placeId is on the map.
	MappedSeats int `json:"mappedSeats"`
}

// InventoryMarket is a schema of the public API.
type InventoryMarket string

// Values of InventoryMarket.
const (
	InventoryMarketPrimary InventoryMarket = "primary"
)

// InventoryChangeKind is a schema of the public API.
type InventoryChangeKind string

// Values of InventoryChangeKind.
const (
	InventoryChangeKindRestock        InventoryChangeKind = "restock"
	InventoryChangeKindDrop           InventoryChangeKind = "drop"
	InventoryChangeKindGaRestock      InventoryChangeKind = "ga_restock"
	InventoryChangeKindGaDrop         InventoryChangeKind = "ga_drop"
	InventoryChangeKindHoldReleased   InventoryChangeKind = "hold_released"
	InventoryChangeKindGaHoldReleased InventoryChangeKind = "ga_hold_released"
	InventoryChangeKindSale           InventoryChangeKind = "sale"
	InventoryChangeKindGaSale         InventoryChangeKind = "ga_sale"
)

// InventoryChangePage is a schema of the public API.
type InventoryChangePage struct {
	Items []InventoryChange `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
}

// InventoryChangeConfidence: `drop`/`ga_drop`: `pending` inside the cart-hold window, `none` for a
// mass transition; `sale`/`ga_sale`: `inferred_high` or `inferred_low`;
// null otherwise.
type InventoryChangeConfidence string

// Values of InventoryChangeConfidence.
const (
	InventoryChangeConfidencePending      InventoryChangeConfidence = "pending"
	InventoryChangeConfidenceInferredHigh InventoryChangeConfidence = "inferred_high"
	InventoryChangeConfidenceInferredLow  InventoryChangeConfidence = "inferred_low"
	InventoryChangeConfidenceNone         InventoryChangeConfidence = "none"
)

// InventoryChange is a schema of the public API.
type InventoryChange struct {
	// Time of the poll that observed the change.
	Ts     time.Time           `json:"ts"`
	Market InventoryMarket     `json:"market"`
	Kind   InventoryChangeKind `json:"kind"`
	// Manifest section name.
	Section string `json:"section"`
	// Manifest section index (-1 when the manifest was not loaded).
	SectionIdx int `json:"sectionIdx"`
	// Row name; empty for GA stock.
	Row string `json:"row"`
	// Seats (or GA units) in the change.
	Count int `json:"count"`
	// Seat labels; empty for GA stock.
	Seats []string `json:"seats"`
	// Ticketmaster placeIds of the seats (join to the seat map's places).
	PlaceIDs []string `json:"placeIds"`
	// List price per ticket in minor units; null when unpriced.
	ListMinor *int64 `json:"listMinor"`
	// All-in price per ticket in minor units; null when unpriced.
	TotalMinor *int64  `json:"totalMinor"`
	Currency   *string `json:"currency"`
	// `drop`/`ga_drop`: `pending` inside the cart-hold window, `none` for a
	// mass transition; `sale`/`ga_sale`: `inferred_high` or `inferred_low`;
	// null otherwise.
	Confidence *InventoryChangeConfidence `json:"confidence"`
	// Part of a mass transition (not counted as sales).
	Mass bool `json:"mass"`
}

// ListingPage is a schema of the public API.
type ListingPage struct {
	Items []Listing `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
	// Time of the latest poll per marketplace, keyed by marketplace (`stubhub`); the listings are as of then.
	AsOf map[string]time.Time `json:"asOf"`
}

// Listing: One listing currently offered on a secondary marketplace. Prices are per
// ticket in minor units of `currency`: `priceMinor` is the seller's price,
// `allInMinor` what a buyer pays with fees (estimated from the marketplace's
// fee model when it does not report it, `allInEstimated`).
type Listing struct {
	// Marketplace (`stubhub`).
	Source string `json:"source"`
	// The marketplace's listing id.
	ExternalID string `json:"externalId"`
	// Section label as the marketplace shows it.
	Section string `json:"section"`
	// Canonical section of the event's seat map; null when unmapped.
	SectionKey *string `json:"sectionKey"`
	// Empty when the marketplace does not tell.
	Row      string  `json:"row"`
	SeatFrom *string `json:"seatFrom"`
	SeatTo   *string `json:"seatTo"`
	// Tickets offered.
	Quantity int `json:"quantity"`
	// Quantities a buyer may purchase; null when any split is allowed or unknown.
	Splits         []int  `json:"splits"`
	PriceMinor     int64  `json:"priceMinor"`
	AllInMinor     int64  `json:"allInMinor"`
	AllInEstimated bool   `json:"allInEstimated"`
	Currency       string `json:"currency"`
	// Listing features the marketplace reports (e.g. aisle, instant delivery).
	Traits      []string  `json:"traits"`
	FirstSeenAt time.Time `json:"firstSeenAt"`
	// Last price or quantity change.
	ChangedAt time.Time `json:"changedAt"`
}

// ListingSectionList is a schema of the public API.
type ListingSectionList struct {
	EventID int64            `json:"eventId"`
	Items   []ListingSection `json:"items"`
	// Time of the latest poll per marketplace, keyed by marketplace (`stubhub`).
	AsOf map[string]time.Time `json:"asOf"`
}

// ListingSection: Current listings of one section across marketplaces. `sources` compares
// the marketplaces offering it, cheapest first; prices are all-in per
// ticket in minor units of each source's currency.
type ListingSection struct {
	// The canonical section key, or the marketplace label of an unmapped section.
	Section string `json:"section"`
	// Canonical section of the event's seat map; null when unmapped.
	SectionKey *string                `json:"sectionKey"`
	Listings   int                    `json:"listings"`
	Tickets    int                    `json:"tickets"`
	Sources    []ListingSectionSource `json:"sources"`
}

// ListingSectionSource is a schema of the public API.
type ListingSectionSource struct {
	// Marketplace (`stubhub`).
	Source   string `json:"source"`
	Listings int    `json:"listings"`
	Tickets  int    `json:"tickets"`
	// Cheapest all-in price per ticket.
	FloorMinor int64  `json:"floorMinor"`
	Currency   string `json:"currency"`
	// Distinct rows listed (empty row labels count as one).
	Rows int `json:"rows"`
}

// ListingChangeKind is a schema of the public API.
type ListingChangeKind string

// Values of ListingChangeKind.
const (
	ListingChangeKindAdded        ListingChangeKind = "added"
	ListingChangeKindPriceChanged ListingChangeKind = "price_changed"
	ListingChangeKindQtyChanged   ListingChangeKind = "qty_changed"
	ListingChangeKindRemoved      ListingChangeKind = "removed"
)

// ListingChangePage is a schema of the public API.
type ListingChangePage struct {
	Items []ListingChange `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
}

// ListingChange: One listing that changed between consecutive polls of a marketplace.
// Prices are all-in per ticket in minor units of `currency`; `qtyOld`/`priceOld`
// are null for `added`, `qtyNew`/`priceNew` null for `removed`.
type ListingChange struct {
	// Time of the poll that observed the change.
	Ts         time.Time         `json:"ts"`
	Source     string            `json:"source"`
	ExternalID string            `json:"externalId"`
	Kind       ListingChangeKind `json:"kind"`
	Section    string            `json:"section"`
	SectionKey *string           `json:"sectionKey"`
	Row        string            `json:"row"`
	QtyOld     *int              `json:"qtyOld"`
	QtyNew     *int              `json:"qtyNew"`
	PriceOld   *int64            `json:"priceOld"`
	PriceNew   *int64            `json:"priceNew"`
	Currency   string            `json:"currency"`
}

// UserRef is a schema of the public API.
type UserRef struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

// Problem: RFC 9457 problem details.
type Problem struct {
	Type     string  `json:"type"`
	Title    string  `json:"title"`
	Status   int     `json:"status"`
	Detail   *string `json:"detail,omitempty"`
	Instance *string `json:"instance,omitempty"`
}

// EventStatus is a schema of the public API.
type EventStatus string

// Values of EventStatus.
const (
	EventStatusOnsale      EventStatus = "onsale"
	EventStatusOffsale     EventStatus = "offsale"
	EventStatusCanceled    EventStatus = "canceled"
	EventStatusPostponed   EventStatus = "postponed"
	EventStatusRescheduled EventStatus = "rescheduled"
	EventStatusUnknown     EventStatus = "unknown"
)

// EventLifecycle: Where the event is in its sale lifecycle, derived when read:
// `archived` 30 days after the start, `past` 2 hours after it (when
// polling stops); cancelled and postponed events stay `announced` until
// then. Otherwise `sold_out` once the sale began and Ticketmaster went
// off sale or the primary inventory is empty, `on_sale` once the public
// on-sale started (or a presale is running, or the event sells without
// an announced on-sale date), `on_sale_soon` when a public on-sale or
// presale starts within 7 days, else `announced`.
type EventLifecycle string

// Values of EventLifecycle.
const (
	EventLifecycleAnnounced  EventLifecycle = "announced"
	EventLifecycleOnSaleSoon EventLifecycle = "on_sale_soon"
	EventLifecycleOnSale     EventLifecycle = "on_sale"
	EventLifecycleSoldOut    EventLifecycle = "sold_out"
	EventLifecyclePast       EventLifecycle = "past"
	EventLifecycleArchived   EventLifecycle = "archived"
)

// EventTier: Polling tier: `watched` (in scope of an enabled market alert rule,
// polled about every 30 s), `standard` (on sale, about hourly) or
// `dormant` (not polled: not on sale yet, ended, cancelled, or without
// a pollable source).
type EventTier string

// Values of EventTier.
const (
	EventTierWatched  EventTier = "watched"
	EventTierStandard EventTier = "standard"
	EventTierDormant  EventTier = "dormant"
)

// Region is a schema of the public API.
type Region struct {
	Code string `json:"code"`
	Name string `json:"name"`
	// Display currency (ISO 4217).
	Currency string `json:"currency"`
	// Venue country codes of the region (Discovery `countryCode` values).
	Countries []string `json:"countries"`
}

// RegionList is a schema of the public API.
type RegionList struct {
	Items []Region `json:"items"`
}

// CatalogFacet is a schema of the public API.
type CatalogFacet struct {
	Name string `json:"name"`
	// Upcoming (or undated) events with this value.
	Events int `json:"events"`
}

// CatalogSegmentFacet is a schema of the public API.
type CatalogSegmentFacet struct {
	CatalogFacet
	Genres []CatalogFacet `json:"genres"`
}

// CatalogFacets is a schema of the public API.
type CatalogFacets struct {
	// Segments by upcoming events, each with its genres by upcoming events.
	Segments []CatalogSegmentFacet `json:"segments"`
}

// EventPage is a schema of the public API.
type EventPage struct {
	Items []EventSummary `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
}

// EventSummary is a schema of the public API.
type EventSummary struct {
	ID int64 `json:"id"`
	// Ticketmaster Discovery event id; empty for events not listed by Ticketmaster Discovery (e.g. `primarySource` `dice`).
	DiscoveryID string `json:"discoveryId"`
	// Ticketmaster host event id (16 uppercase hex); null for events not hosted on Ticketmaster.
	TmEventID *string `json:"tmEventId"`
	// The event's primary seller, whose inventory backs the `primary` market: `tm` when hosted on
	// Ticketmaster, `dice` for a Dice event (see `GET /v1/events/{id}/ticket-types`); null when the
	// platform collects no primary inventory for it.
	PrimarySource *PrimarySource `json:"primarySource"`
	Name          string         `json:"name"`
	URL           *string        `json:"url"`
	// Start instant; null while the date is TBA/TBD.
	StartsAt  *time.Time `json:"startsAt"`
	LocalDate *string    `json:"localDate"`
	// IANA zone of the event.
	Timezone    *string         `json:"timezone"`
	Status      EventStatus     `json:"status"`
	Venue       *VenueRef       `json:"venue,omitempty"`
	Attractions []AttractionRef `json:"attractions"`
	// Classification segment, e.g. Music or Sports.
	Segment         *string     `json:"segment"`
	Genre           *string     `json:"genre"`
	PublicSaleStart *time.Time  `json:"publicSaleStart"`
	PriceRange      *PriceRange `json:"priceRange,omitempty"`
	// When the catalog first saw the event.
	FirstSeenAt time.Time `json:"firstSeenAt"`
	// Set when a complete sweep no longer found the event in Discovery.
	MissingSince *time.Time     `json:"missingSince"`
	Lifecycle    EventLifecycle `json:"lifecycle"`
	Tier         EventTier      `json:"tier"`
	// Whether the event has a compiled seat map.
	HasMap bool `json:"hasMap"`
	// Sources the event is available at, sorted (`tm` when hosted on Ticketmaster, `dice` for a Dice event, then marketplaces with a linked listing).
	AvailableSources []string `json:"availableSources"`
	// Tickets listed on Ticketmaster resale at the last poll of its listings (0: none listed now). Absent
	// for events without Ticketmaster resale data (no linked `tmresale` listing, not polled yet) and for
	// plans without market data.
	TmResaleTickets *int          `json:"tmResaleTickets,omitempty"`
	Metrics         *EventMetrics `json:"metrics,omitempty"`
}

// EventDetail is a schema of the public API.
type EventDetail struct {
	EventSummary
	// Local start time (HH:MM:SS); null when the time is TBA.
	LocalTime *string `json:"localTime"`
	Subgenre  *string `json:"subgenre"`
	// Ticketing system hosting the event (ticketmaster, universe, frontgate, ticketweb, ...).
	Source        *string    `json:"source"`
	Locale        *string    `json:"locale"`
	DateTBA       bool       `json:"dateTBA"`
	DateTBD       bool       `json:"dateTBD"`
	TimeTBA       bool       `json:"timeTBA"`
	PublicSaleEnd *time.Time `json:"publicSaleEnd"`
	PublicSaleTBD bool       `json:"publicSaleTBD"`
	// When Discovery last returned the event.
	LastSeenAt time.Time `json:"lastSeenAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	// Public on-sale and presale windows, ordered by start.
	Sales []SaleWindow `json:"sales"`
	// Secondary marketplaces the event is linked to, with how fresh their listings are.
	Sources []EventSource `json:"sources"`
}

// EventSourceStatus is a schema of the public API.
type EventSourceStatus string

// Values of EventSourceStatus.
const (
	EventSourceStatusLinked EventSourceStatus = "linked"
)

// EventSourceMethod: `exact`: the marketplace reported the event's own id (its TM event id, or the Dice event id of a Dice event); `auto`: matched by score; `manual`: linked by staff.
type EventSourceMethod string

// Values of EventSourceMethod.
const (
	EventSourceMethodExact  EventSourceMethod = "exact"
	EventSourceMethodAuto   EventSourceMethod = "auto"
	EventSourceMethodManual EventSourceMethod = "manual"
)

// EventSourceFreshness: `fresh`: polled within twice its current poll interval; `stale`:
// older (or the event is no longer polled); `pending`: not polled yet.
type EventSourceFreshness string

// Values of EventSourceFreshness.
const (
	EventSourceFreshnessFresh   EventSourceFreshness = "fresh"
	EventSourceFreshnessStale   EventSourceFreshness = "stale"
	EventSourceFreshnessPending EventSourceFreshness = "pending"
)

// EventSource is a schema of the public API.
type EventSource struct {
	// Marketplace source id (`stubhub` — StubHub and viagogo, one catalogue sold under two brands; `crowdvolt` — CrowdVolt, US Dice and Ticketmaster/AXS music events; `tickpick` — TickPick, the whole US catalogue, all-in prices; `vividseats` — Vivid Seats, the whole US catalogue, pre-fee and all-in prices; `ticombo` — Ticombo, a European marketplace that also lists US events, pre-fee and all-in prices in each seller's currency; `tmresale` — Ticketmaster's fan-to-fan resale listings, US, Canadian, UK and Irish events hosted on Ticketmaster, pre-fee and all-in prices).
	Source string `json:"source"`
	// The marketplace's event id.
	ExternalID string `json:"externalId"`
	// The event's page on the marketplace: for `stubhub`, its canonical StubHub
	// link. StubHub and viagogo sell one catalogue under two brands, so the same
	// `externalId` opens on either: `https://www.stubhub.com/event/{externalId}`
	// and `https://www.viagogo.com/E-{externalId}`. For `crowdvolt`,
	// `https://www.crowdvolt.com/event/{externalId}`. For `tickpick`,
	// `https://www.tickpick.com/buy-tickets/{externalId}/`. For `vividseats`,
	// `https://www.vividseats.com/tickets/production/{externalId}`. For `ticombo`,
	// the event page found by the search (no URL form opens an id alone). For `tmresale`,
	// `https://www.ticketmaster.com/event/{externalId}`.
	URL    *string           `json:"url"`
	Status EventSourceStatus `json:"status"`
	// `exact`: the marketplace reported the event's own id (its TM event id, or the Dice event id of a Dice event); `auto`: matched by score; `manual`: linked by staff.
	Method   EventSourceMethod `json:"method"`
	LinkedAt time.Time         `json:"linkedAt"`
	// When the marketplace's listings were last collected; null until the first poll.
	LastPolledAt *time.Time `json:"lastPolledAt"`
	// `fresh`: polled within twice its current poll interval; `stale`:
	// older (or the event is no longer polled); `pending`: not polled yet.
	Freshness EventSourceFreshness `json:"freshness"`
}

// VenueRef is a schema of the public API.
type VenueRef struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	City        *string `json:"city"`
	StateCode   *string `json:"stateCode"`
	CountryCode *string `json:"countryCode"`
}

// AttractionRef is a schema of the public API.
type AttractionRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SearchResults is a schema of the public API.
type SearchResults struct {
	// The folded query the names were matched against.
	Q           string             `json:"q"`
	Events      []SearchEvent      `json:"events"`
	Attractions []SearchAttraction `json:"attractions"`
	Venues      []SearchVenue      `json:"venues"`
}

// SearchEvent is a schema of the public API.
type SearchEvent struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	StartsAt  *time.Time  `json:"startsAt"`
	LocalDate *string     `json:"localDate"`
	Timezone  *string     `json:"timezone"`
	Status    EventStatus `json:"status"`
	Venue     *VenueRef   `json:"venue,omitempty"`
	// Not started yet, or undated.
	Upcoming bool `json:"upcoming"`
}

// SearchAttraction is a schema of the public API.
type SearchAttraction struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Segment        *string `json:"segment"`
	Genre          *string `json:"genre"`
	UpcomingEvents int     `json:"upcomingEvents"`
}

// SearchVenue is a schema of the public API.
type SearchVenue struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	City           *string `json:"city"`
	StateCode      *string `json:"stateCode"`
	CountryCode    *string `json:"countryCode"`
	UpcomingEvents int     `json:"upcomingEvents"`
}

// EventRollup: Catalog events of an attraction or venue around now (upcoming = not started, or undated).
type EventRollup struct {
	Upcoming int `json:"upcoming"`
	Past     int `json:"past"`
	// Start of the first dated upcoming event.
	FirstStartsAt *time.Time `json:"firstStartsAt"`
	// Start of the last dated upcoming event.
	LastStartsAt *time.Time `json:"lastStartsAt"`
	// Local date (in the event's zone) of the first dated upcoming event.
	FirstLocalDate *string `json:"firstLocalDate"`
	// Local date of the last dated upcoming event.
	LastLocalDate *string `json:"lastLocalDate"`
}

// Classification: A Discovery classification; names are empty where Discovery leaves them undefined.
type Classification struct {
	Primary  bool   `json:"primary"`
	Family   bool   `json:"family"`
	Segment  string `json:"segment"`
	Genre    string `json:"genre"`
	Subgenre string `json:"subgenre"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
}

// ExternalLink is a schema of the public API.
type ExternalLink struct {
	// Discovery's link kind, e.g. homepage, spotify, youtube, instagram, wiki.
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

// AttractionHistory is a schema of the public API.
type AttractionHistory struct {
	AttractionID int64                    `json:"attractionId"`
	Items        []AttractionHistoryEvent `json:"items"`
}

// AttractionHistoryEvent: One past event and how its market ended. Spreads and sell-through are
// fractions (0.12 = 12%) with the definitions of `EventMetrics`, as of
// `asOf` (the event's last metrics); null when not measured.
type AttractionHistoryEvent struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	StartsAt time.Time `json:"startsAt"`
	Venue    *VenueRef `json:"venue,omitempty"`
	// Last primary floor, all-in minor units of `currency`.
	PrimaryFloor *int64  `json:"primaryFloor"`
	Currency     *string `json:"currency"`
	// TM resale floor ÷ primary floor − 1.
	ResaleSpread *float64 `json:"resaleSpread"`
	// Best marketplace floor ÷ primary floor − 1; null without secondary marketplaces in the plan.
	SecondarySpread *float64 `json:"secondarySpread"`
	// Primary sell-through, 1 − open primary tickets / seat-map capacity.
	SellThrough *float64 `json:"sellThrough"`
	// (first primary inventory − last primary inventory before the start) ÷
	// first primary inventory, over the collected primary samples; null with
	// fewer than two samples or no first inventory.
	Depletion *float64 `json:"depletion"`
	// First minus last primary inventory over the same samples.
	PrimarySold *int `json:"primarySold"`
	// Time of the newest market sample behind the metrics.
	AsOf *time.Time `json:"asOf"`
}

// AttractionDetail is a schema of the public API.
type AttractionDetail struct {
	ID int64 `json:"id"`
	// Ticketmaster Discovery attraction id.
	DiscoveryID string `json:"discoveryId"`
	Name        string `json:"name"`
	// Ticketmaster attraction page.
	URL             *string          `json:"url"`
	Segment         *string          `json:"segment"`
	Genre           *string          `json:"genre"`
	Subgenre        *string          `json:"subgenre"`
	Classifications []Classification `json:"classifications"`
	// Official and social links, by kind.
	Links  []ExternalLink `json:"links"`
	Rollup EventRollup    `json:"rollup"`
	// Distinct cities of the upcoming events' venues.
	UpcomingCities int `json:"upcomingCities"`
	// Countries (ISO 3166-1 alpha-2) of the upcoming events' venues, alphabetical.
	UpcomingCountries []string      `json:"upcomingCountries"`
	NextEvent         *EventSummary `json:"nextEvent,omitempty"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

// GeoPoint is a schema of the public API.
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// VenueSeatmap: A compiled seat map used by events at the venue.
type VenueSeatmap struct {
	// Hex sha256 of the uncompressed artifact.
	Hash string `json:"hash"`
	// The artifact (`getSeatmap`).
	URL string `json:"url"`
	// The artifact on the public map artifact store, when published there (see `Seatmap.cdnUrl`).
	CdnURL     *string `json:"cdnUrl,omitempty"`
	Format     int     `json:"format"`
	PageWidth  int     `json:"pageWidth"`
	PageHeight int     `json:"pageHeight"`
	Sections   int     `json:"sections"`
	GaSections int     `json:"gaSections"`
	// Seats with a position on the map.
	Places int `json:"places"`
	// Declared places (seats plus GA capacity); null when the map declares none.
	Capacity   *int      `json:"capacity"`
	CompiledAt time.Time `json:"compiledAt"`
	// Events at the venue linked to this map.
	Events int `json:"events"`
	// Latest map fetch of those events.
	LastFetchedAt time.Time `json:"lastFetchedAt"`
}

// VenueConfig: A Ticketmaster venue configuration (layout) seen at the venue.
type VenueConfig struct {
	VenueConfigID  string `json:"venueConfigId"`
	Events         int    `json:"events"`
	UpcomingEvents int    `json:"upcomingEvents"`
	// Compiled maps of this configuration, most recently fetched first.
	Seatmaps []VenueSeatmap `json:"seatmaps"`
}

// VenueDetail is a schema of the public API.
type VenueDetail struct {
	ID int64 `json:"id"`
	// Ticketmaster Discovery venue id; empty for venues not listed by Ticketmaster Discovery (e.g. Dice venues).
	DiscoveryID string `json:"discoveryId"`
	// The seller whose catalog lists the venue, and whose site `url` opens: `tm` for Ticketmaster
	// Discovery venues, `dice` for Dice venues.
	PrimarySource *PrimarySource `json:"primarySource"`
	Name          string         `json:"name"`
	URL           *string        `json:"url"`
	Address       *string        `json:"address"`
	City          *string        `json:"city"`
	StateCode     *string        `json:"stateCode"`
	CountryCode   *string        `json:"countryCode"`
	PostalCode    *string        `json:"postalCode"`
	// IANA zone of the venue.
	Timezone *string   `json:"timezone"`
	Location *GeoPoint `json:"location,omitempty"`
	// Largest declared capacity among the venue's compiled seat maps; null without one.
	Capacity  *int          `json:"capacity"`
	Rollup    EventRollup   `json:"rollup"`
	NextEvent *EventSummary `json:"nextEvent,omitempty"`
	// Venue configurations, most upcoming events first.
	Configs   []VenueConfig `json:"configs"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// PriceRange: Face-value range advertised by Ticketmaster, in minor units (cents).
type PriceRange struct {
	MinMinor int64  `json:"minMinor"`
	MaxMinor int64  `json:"maxMinor"`
	Currency string `json:"currency"`
}

// SaleWindowKind is a schema of the public API.
type SaleWindowKind string

// Values of SaleWindowKind.
const (
	SaleWindowKindPublic  SaleWindowKind = "public"
	SaleWindowKindPresale SaleWindowKind = "presale"
)

// SaleWindow is a schema of the public API.
type SaleWindow struct {
	Kind SaleWindowKind `json:"kind"`
	// Presale name; empty for the public on-sale.
	Name     string     `json:"name"`
	StartsAt *time.Time `json:"startsAt"`
	EndsAt   *time.Time `json:"endsAt"`
	// Presale sign-up or access page, when Discovery provides one.
	URL *string `json:"url"`
	// Discovery's presale terms: who is eligible and whether (and how) a
	// presale code is needed. Discovery does not publish the codes themselves.
	Description *string `json:"description"`
}

// MetricScope: Market scope of a metric: Ticketmaster primary (`primary`), Ticketmaster
// resale (`tm_resale`: Ticketmaster's fan-to-fan resale order book, from
// the `tmresale` listings, with asking and all-in prices and sales inferred
// from listings that sold down or disappeared), one marketplace
// (`secondary:<source>`, e.g. `secondary:stubhub`), or the best of the
// marketplaces (`secondary:best`: lowest floor, summed inventory and sales
// across them, in the event's currency; TM resale is not part of it).
// Secondary scopes need the secondary-marketplaces feature. More
// marketplaces may be added; treat unknown values as other scopes.
type MetricScope string

// Values of MetricScope.
const (
	MetricScopePrimary             MetricScope = "primary"
	MetricScopeTmResale            MetricScope = "tm_resale"
	MetricScopeSecondaryBest       MetricScope = "secondary:best"
	MetricScopeSecondaryStubhub    MetricScope = "secondary:stubhub"
	MetricScopeSecondaryCrowdvolt  MetricScope = "secondary:crowdvolt"
	MetricScopeSecondaryTickpick   MetricScope = "secondary:tickpick"
	MetricScopeSecondaryVividseats MetricScope = "secondary:vividseats"
	MetricScopeSecondaryTicombo    MetricScope = "secondary:ticombo"
)

// EventMetricsMovementDir: Sign of chg24h.
type EventMetricsMovementDir int

// Values of EventMetricsMovementDir.
const (
	EventMetricsMovementDir1Desc EventMetricsMovementDir = -1
	EventMetricsMovementDir0     EventMetricsMovementDir = 0
	EventMetricsMovementDir1     EventMetricsMovementDir = 1
)

// EventMetrics: Serving metrics of one event and scope, recomputed whenever the
// event's market data changes and on every sweep. Money is minor units of
// `currency`, all-in per ticket; changes, depletion, sell-through, churn and
// spread are fractions (0.12 = 12%). Sales count settled sales of every
// confidence (`confirmed`, `inferred_high` and `inferred_low`), mass
// transitions excluded. `secondary:<source>` is in the marketplace's
// currency; `secondary:best` converts each marketplace's prices into the
// event's currency (its Ticketmaster currency) at the daily reference rate of
// each value's day, leaving out marketplaces without a rate.
type EventMetrics struct {
	Scope MetricScope `json:"scope"`
	// Movement cohort, segment slug and days-to-event bucket (`music/3-7d`; buckets `0-2d`, `3-7d`, `8-30d`, `31-90d`, `90d+`, `past`, `tba`).
	Cohort   string  `json:"cohort"`
	Currency *string `json:"currency"`
	// Lowest all-in price of the latest sample; null without priced stock. `secondary:best`, the lowest across marketplaces.
	Floor *int64 `json:"floor"`
	// Lowest all-in price among stock that can sell two tickets together:
	// resale and marketplace listings that can split to a pair; on
	// Ticketmaster primary the cheapest price level with two adjacent open
	// seats of one row (seat ids consecutive in the manifest) or a GA
	// section with two open places (`secondary:best`, the lowest across
	// marketplaces).
	FloorPair *int64 `json:"floorPair"`
	// Ticket-weighted 25th percentile all-in price of the latest sample.
	// `secondary:best`: the ticket-weighted median of the marketplaces' p25s.
	P25 *int64 `json:"p25"`
	// Ticket-weighted median all-in price of the latest sample. `secondary:best`:
	// the ticket-weighted median of the marketplaces' medians.
	Median *int64 `json:"median"`
	// Tickets available in the latest sample (primary, open seats plus GA stock).
	Inventory *int `json:"inventory"`
	// Time-weighted mean floor over 24h (each 5-minute sample holds until the next; gaps carry forward).
	Avg24h *int64 `json:"avg24h"`
	Avg72h *int64 `json:"avg72h"`
	Avg7d  *int64 `json:"avg7d"`
	Avg30d *int64 `json:"avg30d"`
	// (floor now − floor 24h ago) / floor 24h ago; null without a baseline.
	Chg24h *float64 `json:"chg24h"`
	Chg72h *float64 `json:"chg72h"`
	Chg7d  *float64 `json:"chg7d"`
	// Tickets sold in the last 24h.
	Sales24h int `json:"sales24h"`
	Sales72h int `json:"sales72h"`
	// Σ qty·all-in of the last 24h's priced sales.
	Volume24h int64 `json:"volume24h"`
	// Σ qty·all-in of the last 72h's priced sales.
	Volume72h int64 `json:"volume72h"`
	// Quantity-weighted mean all-in price of the last 24h's priced sales.
	SaleAvg24h *int64 `json:"saleAvg24h"`
	SaleAvg72h *int64 `json:"saleAvg72h"`
	// Tickets sold per hour, exponentially weighted with a 6h time constant.
	Velocity float64 `json:"velocity"`
	// (inventory 24h ago − inventory now) / inventory 24h ago.
	Depletion24h *float64 `json:"depletion24h"`
	// Primary only, 1 − open primary tickets / seat-map capacity; null without a compiled seat map.
	SellThrough *float64 `json:"sellThrough"`
	// Primary only: hours until the open primary tickets sell out at the current pace,
	// `inventory` ÷ `velocity` (settled sales; Ticketmaster settles a sale ~20 minutes after the
	// seats leave the inventory into a cart, so in a rush the estimate errs late). Null when there
	// are no open tickets (sold out or not on sale), without a settled sale in the last 24h, when
	// the latest primary sample is over 2h old, and when the estimate is past 90 days or the
	// event start.
	SelloutHours *float64 `json:"selloutHours"`
	// Inventory turnover over 24h, (tickets added + tickets gone for good) / inventory.
	Churn24h *float64 `json:"churn24h"`
	// floor ÷ primary floor − 1 (the primary floor converted into this scope's
	// currency at the day's reference rate); null for `primary` and without
	// both floors.
	Spread *float64 `json:"spread"`
	// 0.35·|z(chg24h)| + 0.30·z(sales24h) + 0.20·z(depletion24h) + 0.15·z(churn24h), z-scores
	// within the cohort (the whole scope for cohorts under 8 events); weights are configurable.
	Movement *float64 `json:"movement"`
	// Sign of chg24h.
	MovementDir EventMetricsMovementDir `json:"movementDir"`
	// Floor at 48 instants 90 minutes apart over the last 72h, oldest first; null where unknown.
	Spark []*float64 `json:"spark"`
	// Time of the newest market sample the metrics are computed from.
	AsOf time.Time `json:"asOf"`
	// Timeframe windows with their comparisons (the 7-day sales, volume and
	// mean sale price are `windows.7d`); null until the row is recomputed
	// after their introduction.
	Windows *MetricWindows `json:"windows"`
}

// MetricWindows: The metrics of each timeframe ending now, keyed by window. Each window is
// compared with the next longer one (`6h` with `12h`, `12h` with `24h`,
// `24h` with `7d`, `7d` with `30d`) and `30d` with the 30 days before it
// (`prev30d`).
type MetricWindows struct {
	N6h  MetricWindow `json:"6h"`
	N12h MetricWindow `json:"12h"`
	N24h MetricWindow `json:"24h"`
	N7d  MetricWindow `json:"7d"`
	N30d MetricWindow `json:"30d"`
}

// MetricWindow: One timeframe's metrics over (now − window, now], defined like the 24h
// fields of EventMetrics (`avg` like `avg24h`, `chg` like `chg24h`, `sales`
// like `sales24h`, …). Money is minor units of the row's `currency`.
type MetricWindow struct {
	// Time-weighted mean floor over the window.
	Avg *int64 `json:"avg"`
	// (floor now − floor at the window start) / floor at the window start; null without a baseline.
	Chg *float64 `json:"chg"`
	// Tickets sold in the window.
	Sales int `json:"sales"`
	// Σ qty·all-in of the window's priced sales.
	Volume int64 `json:"volume"`
	// Quantity-weighted mean all-in price of the window's priced sales.
	SaleAvg *int64 `json:"saleAvg"`
	// sales ÷ the window's covered hours (those since the scope's first market sample or sale; 0 without any).
	SalesPerHour float64 `json:"salesPerHour"`
	// (inventory at the window start − inventory now) / inventory at the window start.
	Depletion *float64               `json:"depletion"`
	Vs        MetricWindowComparison `json:"vs"`
}

// MetricWindowComparisonWindow: The base, the next longer window ending now; `prev30d` is the 30 days before the `30d` window.
type MetricWindowComparisonWindow string

// Values of MetricWindowComparisonWindow.
const (
	MetricWindowComparisonWindow12h     MetricWindowComparisonWindow = "12h"
	MetricWindowComparisonWindow24h     MetricWindowComparisonWindow = "24h"
	MetricWindowComparisonWindow7d      MetricWindowComparisonWindow = "7d"
	MetricWindowComparisonWindow30d     MetricWindowComparisonWindow = "30d"
	MetricWindowComparisonWindowPrev30d MetricWindowComparisonWindow = "prev30d"
)

// MetricWindowComparison: A window against its base (`window`): fractions (0.12 = 12% higher), price
// means as ratios and sales and volume as per-hour rates, each value ÷ the
// base's − 1; null when either side is unknown or the base is zero. Rates
// count only the hours since the scope's first market sample or sale, so a
// recently tracked event is not compared against hours it had no data for.
type MetricWindowComparison struct {
	// The base, the next longer window ending now; `prev30d` is the 30 days before the `30d` window.
	Window MetricWindowComparisonWindow `json:"window"`
	// avg ÷ the base's avg − 1.
	Avg *float64 `json:"avg"`
	// (sales ÷ covered hours) ÷ (base sales ÷ base covered hours) − 1; null when the base sold nothing or has no covered hours.
	SalesRate *float64 `json:"salesRate"`
	// (volume ÷ covered hours) ÷ (base volume ÷ base covered hours) − 1; null when the base has no volume or no covered hours.
	VolumeRate *float64 `json:"volumeRate"`
	// saleAvg ÷ the base's saleAvg − 1.
	SaleAvg *float64 `json:"saleAvg"`
}

// EventMetricsList is a schema of the public API.
type EventMetricsList struct {
	EventID int64 `json:"eventId"`
	// One entry per scope with market data.
	Items []EventMetrics `json:"items"`
}

// MoverPageWindow is a schema of the public API.
type MoverPageWindow string

// Values of MoverPageWindow.
const (
	MoverPageWindow24h MoverPageWindow = "24h"
	MoverPageWindow72h MoverPageWindow = "72h"
	MoverPageWindow7d  MoverPageWindow = "7d"
)

// MoverPage is a schema of the public API.
type MoverPage struct {
	Window MoverPageWindow `json:"window"`
	Scope  MetricScope     `json:"scope"`
	Items  []Mover         `json:"items"`
}

// Mover is a schema of the public API.
type Mover struct {
	Event EventSummary `json:"event"`
	// 24h, the movement score; 72h/7d, |z| of the window's floor change within the cohort.
	Score float64 `json:"score"`
	// Floor change over the window (fraction).
	Change *float64 `json:"change"`
}

// SeriesResolution is a schema of the public API.
type SeriesResolution string

// Values of SeriesResolution.
const (
	SeriesResolution5m SeriesResolution = "5m"
	SeriesResolution1h SeriesResolution = "1h"
	SeriesResolution1d SeriesResolution = "1d"
)

// SeriesMetric is a schema of the public API.
type SeriesMetric string

// Values of SeriesMetric.
const (
	SeriesMetricFloor    SeriesMetric = "floor"
	SeriesMetricFloorMin SeriesMetric = "floorMin"
	SeriesMetricMedian   SeriesMetric = "median"
	SeriesMetricTickets  SeriesMetric = "tickets"
	SeriesMetricSales    SeriesMetric = "sales"
	SeriesMetricVolume   SeriesMetric = "volume"
)

// EventSeries: Column-oriented chart data: `t` holds the bucket starts and each requested
// metric one value per bucket (null where the bucket has no data). `floor`
// is the bucket's last floor, `floorMin` its lowest, `median` the mean
// median, `tickets` the last inventory, `sales` the tickets sold and
// `volume` Σ qty·all-in of its priced sales (all-in minor units).
type EventSeries struct {
	EventID int64       `json:"eventId"`
	Scope   MetricScope `json:"scope"`
	// Currency of the money columns (`floor`, `floorMin`, `median`,
	// `volume`): the scope's latest snapshot currency as of `to` (the
	// marketplace's own for `secondary:<source>`; the event's for
	// `secondary:best`, which converts). Null when the range has no buckets.
	Currency *string          `json:"currency"`
	Res      SeriesResolution `json:"res"`
	From     time.Time        `json:"from"`
	To       time.Time        `json:"to"`
	T        []time.Time      `json:"t"`
	Floor    *SeriesColumn    `json:"floor,omitempty"`
	FloorMin *SeriesColumn    `json:"floorMin,omitempty"`
	Median   *SeriesColumn    `json:"median,omitempty"`
	Tickets  *SeriesColumn    `json:"tickets,omitempty"`
	Sales    *SeriesColumn    `json:"sales,omitempty"`
	Volume   *SeriesColumn    `json:"volume,omitempty"`
}

// SeriesColumn is a schema of the public API.
type SeriesColumn []*int64

// SaleConfidence: `confirmed`: reported by the marketplace; `inferred_high`/`inferred_low`:
// inferred from consecutive polls (low after a long gap between them).
type SaleConfidence string

// Values of SaleConfidence.
const (
	SaleConfidenceConfirmed    SaleConfidence = "confirmed"
	SaleConfidenceInferredHigh SaleConfidence = "inferred_high"
	SaleConfidenceInferredLow  SaleConfidence = "inferred_low"
)

// SaleKind is a schema of the public API.
type SaleKind string

// Values of SaleKind.
const (
	SaleKindSale      SaleKind = "sale"
	SaleKindGaSale    SaleKind = "ga_sale"
	SaleKindConfirmed SaleKind = "confirmed"
	SaleKindPartial   SaleKind = "partial"
	SaleKindRemoval   SaleKind = "removal"
)

// Sale: One settled sale: Ticketmaster primary seats (`sale`) or GA units (`ga_sale`)
// that left `open` and stayed gone past the cart-hold window, or a resale
// sale on Ticketmaster resale (scope `tm_resale`) or a marketplace (scope
// `secondary:<source>`): reported by the source (`confirmed`), a listing
// whose quantity dropped (`partial`) or one that vanished and stayed gone
// past the grace period (`removal`). Prices are per ticket; `priceMinor`
// is all-in.
type Sale struct {
	// Time of the poll that settled the sale.
	Ts      time.Time   `json:"ts"`
	EventID int64       `json:"eventId"`
	Scope   MetricScope `json:"scope"`
	Kind    SaleKind    `json:"kind"`
	Section string      `json:"section"`
	// Empty for GA stock.
	Row   string   `json:"row"`
	Seats []string `json:"seats"`
	Qty   int      `json:"qty"`
	// All-in price per ticket; null when unpriced. Marketplace sales without a reported all-in price estimate it from the source's buyer fee (`allInEstimated`).
	PriceMinor *int64  `json:"priceMinor"`
	ListMinor  *int64  `json:"listMinor"`
	Currency   *string `json:"currency"`
	// `priceMinor` adds an estimated buyer fee to the seller's price (`listMinor`); false for Ticketmaster sales.
	AllInEstimated bool           `json:"allInEstimated"`
	Confidence     SaleConfidence `json:"confidence"`
	Event          *SaleEventRef  `json:"event,omitempty"`
}

// SaleEventRef: Summary of the sale's event (cross-event feeds).
type SaleEventRef struct {
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	StartsAt *time.Time `json:"startsAt"`
	Segment  *string    `json:"segment"`
	Genre    *string    `json:"genre"`
	Venue    *string    `json:"venue"`
	City     *string    `json:"city"`
	Country  *string    `json:"country"`
}

// SalePage is a schema of the public API.
type SalePage struct {
	Items []Sale `json:"items"`
	// Cursor of the next page; null on the last page.
	NextCursor *string `json:"nextCursor"`
}

// SectionSalesList is a schema of the public API.
type SectionSalesList struct {
	EventID int64 `json:"eventId"`
	// One row per section and currency, most tickets first.
	Items []SectionSales `json:"items"`
}

// SectionSales: Settled sales of one section in one currency. Money is minor units of
// `currency`, all-in per ticket; prices aggregate the priced sales only
// (`pricedTickets`).
type SectionSales struct {
	// Manifest section name (Ticketmaster) or canonical section key (marketplaces; their label when unmapped).
	Section string `json:"section"`
	// Null when none of the section's sales is priced.
	Currency *string `json:"currency"`
	// Sale rows (a seat group or a marketplace sale each).
	Sales         int `json:"sales"`
	Tickets       int `json:"tickets"`
	PricedTickets int `json:"pricedTickets"`
	// Σ qty·price of the priced sales.
	VolumeMinor int64 `json:"volumeMinor"`
	// Ticket-weighted mean price (`volumeMinor` ÷ `pricedTickets`, rounded).
	AvgMinor   *int64    `json:"avgMinor"`
	MinMinor   *int64    `json:"minMinor"`
	MaxMinor   *int64    `json:"maxMinor"`
	LastSaleAt time.Time `json:"lastSaleAt"`
	// Tickets per scope (sale scopes only, most first).
	Scopes []SectionSalesScope `json:"scopes"`
}

// SectionSalesScope is a schema of the public API.
type SectionSalesScope struct {
	Scope   MetricScope `json:"scope"`
	Tickets int         `json:"tickets"`
}

// SalesSummaryWindow is a schema of the public API.
type SalesSummaryWindow string

// Values of SalesSummaryWindow.
const (
	SalesSummaryWindow1h  SalesSummaryWindow = "1h"
	SalesSummaryWindow6h  SalesSummaryWindow = "6h"
	SalesSummaryWindow24h SalesSummaryWindow = "24h"
	SalesSummaryWindow72h SalesSummaryWindow = "72h"
	SalesSummaryWindow7d  SalesSummaryWindow = "7d"
)

// SalesSummary is a schema of the public API.
type SalesSummary struct {
	Window SalesSummaryWindow `json:"window"`
	// Start of the window (inclusive).
	From time.Time `json:"from"`
	// End of the window (exclusive), the time of the request.
	To time.Time `json:"to"`
	// One entry per scope (filtered by `scope` when given).
	Markets []SalesMarketTotal `json:"markets"`
}

// SalesMarketTotal is a schema of the public API.
type SalesMarketTotal struct {
	Scope MetricScope `json:"scope"`
	// Whether the scope's sales are collected in the requested `country` or `region` (and `segment`) at all. False when no source covers them, e.g. `tm_resale` outside the US, Canada, the UK and Ireland: its zeros then mean not collected, not no sales. Always true without a country or region.
	Tracked bool `json:"tracked"`
	// Settled sales (rows of the feed).
	Sales int `json:"sales"`
	// Σ qty.
	Tickets int64 `json:"tickets"`
	// Σ qty·all-in of the priced sales per currency, largest first.
	Volume []MoneyTotal `json:"volume"`
}

// MoneyTotal is a schema of the public API.
type MoneyTotal struct {
	Currency string `json:"currency"`
	Minor    int64  `json:"minor"`
}

// CompareSeriesSetAlign is a schema of the public API.
type CompareSeriesSetAlign string

// Values of CompareSeriesSetAlign.
const (
	CompareSeriesSetAlignCalendar    CompareSeriesSetAlign = "calendar"
	CompareSeriesSetAlignDaysToEvent CompareSeriesSetAlign = "days_to_event"
)

// CompareSeriesSet is a schema of the public API.
type CompareSeriesSet struct {
	Metric SeriesMetric          `json:"metric"`
	Align  CompareSeriesSetAlign `json:"align"`
	Scope  MetricScope           `json:"scope"`
	Res    SeriesResolution      `json:"res"`
	Series []CompareSeries       `json:"series"`
}

// CompareSeries: One event's series. `t` are bucket starts, `daysToEvent` the same instants
// as days relative to the event's start (negative before the show; null
// without a start date), `v` the metric's values.
type CompareSeries struct {
	Event       SaleEventRef `json:"event"`
	T           []time.Time  `json:"t"`
	DaysToEvent []*float64   `json:"daysToEvent"`
	V           SeriesColumn `json:"v"`
}

// PredictionStatus is a schema of the public API.
type PredictionStatus string

// Values of PredictionStatus.
const (
	PredictionStatusOpen PredictionStatus = "open"
	PredictionStatusHit  PredictionStatus = "hit"
	PredictionStatusMiss PredictionStatus = "miss"
	PredictionStatusVoid PredictionStatus = "void"
)

// PredictionStatusFilter: `graded`: hit or miss.
type PredictionStatusFilter string

// Values of PredictionStatusFilter.
const (
	PredictionStatusFilterOpen   PredictionStatusFilter = "open"
	PredictionStatusFilterGraded PredictionStatusFilter = "graded"
	PredictionStatusFilterHit    PredictionStatusFilter = "hit"
	PredictionStatusFilterMiss   PredictionStatusFilter = "miss"
	PredictionStatusFilterVoid   PredictionStatusFilter = "void"
)

// PredictionCategory: The call family, each with its own model versions, track record and publication: `tm` (Ticketmaster primary, read with its ticket counts) or `dice` (DICE primary resold on a marketplace, scored mostly on the lineup's artists and the venue).
type PredictionCategory string

// Values of PredictionCategory.
const (
	PredictionCategoryTm   PredictionCategory = "tm"
	PredictionCategoryDice PredictionCategory = "dice"
)

// PredictionEvent is a schema of the public API.
type PredictionEvent struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	StartsAt  *time.Time `json:"startsAt,omitempty"`
	Status    string     `json:"status"`
	VenueName *string    `json:"venueName,omitempty"`
	City      *string    `json:"city,omitempty"`
}

// PredictionCall: A scalping call as customers see it (the model's inputs left out).
type PredictionCall struct {
	ID       int64              `json:"id"`
	Category PredictionCategory `json:"category"`
	Event    PredictionEvent    `json:"event"`
	// Version of the logic that made the call (graded against its own required markup and risk at the call).
	Model  string           `json:"model"`
	Status PredictionStatus `json:"status"`
	// Latest chance (%) that the call is wrong; frozen when the event starts.
	Risk int `json:"risk"`
	// The risk when the call was made; graded for calibration.
	RiskAtCall int `json:"riskAtCall"`
	// Why, strongest factor first (the risks included).
	Reasons  []string `json:"reasons"`
	Currency string   `json:"currency"`
	// Primary floor when the call was made (minor units), the price the call buys at.
	Entry int64 `json:"entry"`
	// Latest resale floor, converted into the currency.
	Resale *int64 `json:"resale,omitempty"`
	// Latest resale ÷ entry − 1.
	Markup *float64 `json:"markup,omitempty"`
	// Markup the resale has to clear for a hit (fees and margin).
	RequiredMarkup float64    `json:"requiredMarkup"`
	CalledAt       time.Time  `json:"calledAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	GradedAt       *time.Time `json:"gradedAt,omitempty"`
	// Median all-in price of the resale sales (marketplaces and Ticketmaster resale) between the call and the event, in the currency. Absent on a void (fewer than 3 resale sales).
	Exit *int64 `json:"exit,omitempty"`
	// exit ÷ entry − 1.
	ExitMarkup *float64 `json:"exitMarkup,omitempty"`
	// How the call was graded, in words.
	Outcome *string `json:"outcome,omitempty"`
}

// PredictionCallPage is a schema of the public API.
type PredictionCallPage struct {
	Items []PredictionCall `json:"items"`
	// Pass as `cursor` for the next page; absent on the last page.
	NextCursor *string `json:"nextCursor,omitempty"`
}

// PredictionBand: Graded calls of one band of the risk at the call.
type PredictionBand struct {
	MinRisk int `json:"minRisk"`
	// Exclusive, except 100.
	MaxRisk int `json:"maxRisk"`
	// Hits and misses.
	Graded int `json:"graded"`
	Misses int `json:"misses"`
	// Mean risk at the call ÷ 100, the miss rate the risk predicted.
	ExpectedMissRate float64 `json:"expectedMissRate"`
	// misses ÷ graded, the miss rate observed.
	MissRate float64 `json:"missRate"`
}

// PredictionRecord: How a set of calls turned out. Graded calls are hits and misses; voids count apart.
type PredictionRecord struct {
	Hits   int `json:"hits"`
	Misses int `json:"misses"`
	Voided int `json:"voided"`
	// hits ÷ (hits + misses); absent before the first grade.
	HitRate *float64 `json:"hitRate,omitempty"`
	// Mean exit markup (realized resale ÷ entry − 1) of the graded calls.
	MeanExitMarkup *float64 `json:"meanExitMarkup,omitempty"`
	// Mean squared error of 1 − risk at the call ÷ 100 against the outcome (0 is perfect, 0.25 is a coin flip).
	Brier *float64 `json:"brier,omitempty"`
	// Calibration per 10-point band of the risk at the call (90–100 includes 100); empty bands left out.
	Bands []PredictionBand `json:"bands"`
}

// PredictionScope: The track record of one model version, or of every version.
type PredictionScope struct {
	// The model version; absent for the overall record.
	Model *string `json:"model,omitempty"`
	// Open calls.
	Open       int              `json:"open"`
	AllTime    PredictionRecord `json:"allTime"`
	Last30Days PredictionRecord `json:"last30Days"`
	Last90Days PredictionRecord `json:"last90Days"`
}

// PredictionTrackRecord is a schema of the public API.
type PredictionTrackRecord struct {
	ComputedAt time.Time          `json:"computedAt"`
	Category   PredictionCategory `json:"category"`
	// The model version making the category's new calls.
	CurrentModel string          `json:"currentModel"`
	Overall      PredictionScope `json:"overall"`
	// Per model version, newest call first.
	Models []PredictionScope `json:"models"`
}

// PredictionAvailabilityReason: Why not: `unpublished` (not published to customers yet); `plan` (the plan lacks `predictions`) no longer occurs, every plan has them.
type PredictionAvailabilityReason string

// Values of PredictionAvailabilityReason.
const (
	PredictionAvailabilityReasonPlan        PredictionAvailabilityReason = "plan"
	PredictionAvailabilityReasonUnpublished PredictionAvailabilityReason = "unpublished"
)

// PredictionCategoryAvailability: Whether the caller sees one category's calls.
type PredictionCategoryAvailability struct {
	Category     PredictionCategory            `json:"category"`
	Available    bool                          `json:"available"`
	Reason       *PredictionAvailabilityReason `json:"reason,omitempty"`
	CurrentModel string                        `json:"currentModel"`
	// Graded calls (hits and misses) of the category's current model.
	Graded int `json:"graded"`
	// Graded calls the current model needs before customers see the category.
	MinGraded int `json:"minGraded"`
}

// PredictionAvailability: The caller's access. The top-level fields are the default category's (`tm`); `categories` has every category, each published on its own.
type PredictionAvailability struct {
	// Whether the caller sees the `tm` predictions.
	Available    bool                          `json:"available"`
	Reason       *PredictionAvailabilityReason `json:"reason,omitempty"`
	CurrentModel string                        `json:"currentModel"`
	// Graded calls (hits and misses) of the current `tm` model.
	Graded int `json:"graded"`
	// Graded calls the current model needs before customers see the predictions.
	MinGraded int `json:"minGraded"`
	// Every category, `tm` first.
	Categories []PredictionCategoryAvailability `json:"categories"`
}

// MustCopStatus: `open` until graded; `hit` sold out within two hours of its public on-sale, `miss` did not, `void` the on-sale went unobserved or the event was canceled, postponed or delisted.
type MustCopStatus string

// Values of MustCopStatus.
const (
	MustCopStatusOpen MustCopStatus = "open"
	MustCopStatusHit  MustCopStatus = "hit"
	MustCopStatusMiss MustCopStatus = "miss"
	MustCopStatusVoid MustCopStatus = "void"
)

// MustCopStatusFilter: `graded`: hit or miss.
type MustCopStatusFilter string

// Values of MustCopStatusFilter.
const (
	MustCopStatusFilterOpen   MustCopStatusFilter = "open"
	MustCopStatusFilterGraded MustCopStatusFilter = "graded"
	MustCopStatusFilterHit    MustCopStatusFilter = "hit"
	MustCopStatusFilterMiss   MustCopStatusFilter = "miss"
	MustCopStatusFilterVoid   MustCopStatusFilter = "void"
)

// MustCopRelease: A past public on-sale of the act with a known outcome.
type MustCopRelease struct {
	EventID   int64      `json:"eventId"`
	Name      string     `json:"name"`
	StartsAt  *time.Time `json:"startsAt,omitempty"`
	VenueName *string    `json:"venueName,omitempty"`
	City      *string    `json:"city,omitempty"`
	// The public on-sale start.
	OnSaleAt time.Time `json:"onSaleAt"`
	// When the primary market was seen sold out, at most two hours after `onSaleAt` (before it: sold out in a presale); absent when it was still on sale then. Unwatched events are checked hourly, so this is an upper bound.
	SoldOutAt *time.Time `json:"soldOutAt,omitempty"`
}

// MustCop: An upcoming public on-sale of an act whose latest known on-sales (Ticketmaster or DICE) mostly sold out instantly: at least 2 of its last 3 within two hours.
type MustCop struct {
	ID    int64           `json:"id"`
	Event PredictionEvent `json:"event"`
	// The headliner the pick rests on.
	Act string `json:"act"`
	// The event's public on-sale start (followed while the pick is open).
	OnSaleAt time.Time     `json:"onSaleAt"`
	Status   MustCopStatus `json:"status"`
	Reasons  []string      `json:"reasons"`
	// The act's on-sales the pick rests on, newest first, as they were at the pick.
	Releases  []MustCopRelease `json:"releases"`
	CalledAt  time.Time        `json:"calledAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	GradedAt  *time.Time       `json:"gradedAt,omitempty"`
	// A hit's sell-out (see `MustCopRelease.soldOutAt`).
	SoldOutAt *time.Time `json:"soldOutAt,omitempty"`
}

// MustCopRecord: How the graded picks turned out. Graded picks are hits and misses; voids count apart.
type MustCopRecord struct {
	Hits   int `json:"hits"`
	Misses int `json:"misses"`
	Voided int `json:"voided"`
	// hits ÷ (hits + misses); absent before the first grade.
	HitRate *float64 `json:"hitRate,omitempty"`
}

// MustCopPage is a schema of the public API.
type MustCopPage struct {
	Items []MustCop `json:"items"`
	// Pass as `cursor` for the next page; absent on the last page.
	NextCursor *string       `json:"nextCursor,omitempty"`
	Record     MustCopRecord `json:"record"`
}

// SavedListRef is a schema of the public API.
type SavedListRef struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ItemCount int    `json:"itemCount"`
}

// SavedList: An org's saved list of events (spreadsheet source) and its saved view.
type SavedList struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// The creator; null when their account is gone.
	Owner     *UserRef      `json:"owner"`
	ItemCount int           `json:"itemCount"`
	View      SavedListView `json:"view"`
	// Whether the caller may delete it (creator or org admin).
	CanDelete bool      `json:"canDelete"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SavedListList is a schema of the public API.
type SavedListList struct {
	Items []SavedList `json:"items"`
}

// SavedListView: Spreadsheet state shared by the org's members.
type SavedListView struct {
	// Visible columns in order (empty = defaults); the event name column is always first.
	Columns []SavedListColumn  `json:"columns"`
	Sort    []SavedListSortKey `json:"sort"`
	Filters []SavedListFilter  `json:"filters"`
	// Conditional formatting rules in priority order.
	Formats []SavedListFormat `json:"formats"`
	GroupBy SavedListGroupBy  `json:"groupBy"`
	Scope   MetricScope       `json:"scope"`
}

// SavedListViewPatch: Present fields replace the saved ones.
type SavedListViewPatch struct {
	Columns []SavedListColumn  `json:"columns,omitempty"`
	Sort    []SavedListSortKey `json:"sort,omitempty"`
	Filters []SavedListFilter  `json:"filters,omitempty"`
	Formats []SavedListFormat  `json:"formats,omitempty"`
	GroupBy *SavedListGroupBy  `json:"groupBy,omitempty"`
	Scope   *MetricScope       `json:"scope,omitempty"`
}

// SavedListGroupBy is a schema of the public API.
type SavedListGroupBy string

// Values of SavedListGroupBy.
const (
	SavedListGroupByNone       SavedListGroupBy = "none"
	SavedListGroupByAttraction SavedListGroupBy = "attraction"
	SavedListGroupByDate       SavedListGroupBy = "date"
	SavedListGroupByVenue      SavedListGroupBy = "venue"
)

// SavedListColumn is a schema of the public API.
type SavedListColumn struct {
	ID string `json:"id"`
	// Pixels; absent = automatic.
	Width *int `json:"width,omitempty"`
}

// SavedListSortKey is a schema of the public API.
type SavedListSortKey struct {
	Column string `json:"column"`
	Desc   bool   `json:"desc"`
}

// SavedListFilter is a schema of the public API.
type SavedListFilter struct {
	Column string            `json:"column"`
	Op     SavedListFilterOp `json:"op"`
	// Number in API units for numeric columns, text otherwise; absent for `empty`/`notEmpty`.
	Value *string `json:"value,omitempty"`
}

// SavedListFilterOp is a schema of the public API.
type SavedListFilterOp string

// Values of SavedListFilterOp.
const (
	SavedListFilterOpLt          SavedListFilterOp = "lt"
	SavedListFilterOpLte         SavedListFilterOp = "lte"
	SavedListFilterOpGt          SavedListFilterOp = "gt"
	SavedListFilterOpGte         SavedListFilterOp = "gte"
	SavedListFilterOpEq          SavedListFilterOp = "eq"
	SavedListFilterOpNe          SavedListFilterOp = "ne"
	SavedListFilterOpContains    SavedListFilterOp = "contains"
	SavedListFilterOpNotContains SavedListFilterOp = "notContains"
	SavedListFilterOpEmpty       SavedListFilterOp = "empty"
	SavedListFilterOpNotEmpty    SavedListFilterOp = "notEmpty"
)

// SavedListFormatColor: Design token the cell is tinted with.
type SavedListFormatColor string

// Values of SavedListFormatColor.
const (
	SavedListFormatColorUp     SavedListFormatColor = "up"
	SavedListFormatColorDown   SavedListFormatColor = "down"
	SavedListFormatColorWarn   SavedListFormatColor = "warn"
	SavedListFormatColorAccent SavedListFormatColor = "accent"
)

// SavedListFormat: Conditional formatting rule: cells of `column` whose value satisfies `op`/`value` (the filter
// semantics) are tinted with `color`. Rules are in priority order; the first match colours a cell.
type SavedListFormat struct {
	Column string            `json:"column"`
	Op     SavedListFilterOp `json:"op"`
	// Number in API units for numeric columns, text otherwise; absent for `empty`/`notEmpty`.
	Value *string `json:"value,omitempty"`
	// Design token the cell is tinted with.
	Color SavedListFormatColor `json:"color"`
}

// CreateSavedListRequest is a schema of the public API.
type CreateSavedListRequest struct {
	// 1–80 characters.
	Name string `json:"name"`
	// At most 500 characters.
	Description *string             `json:"description,omitempty"`
	EventIDs    []int64             `json:"eventIds,omitempty"`
	View        *SavedListViewPatch `json:"view,omitempty"`
}

// UpdateSavedListRequest is a schema of the public API.
type UpdateSavedListRequest struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	View        *SavedListViewPatch `json:"view,omitempty"`
}

// SavedListItem is a schema of the public API.
type SavedListItem struct {
	EventID int64  `json:"eventId"`
	Note    string `json:"note"`
	// Cost per ticket in minor units of `currency`.
	CostBasis *int64    `json:"costBasis"`
	Currency  *string   `json:"currency"`
	QtyHeld   *int      `json:"qtyHeld"`
	AddedBy   *UserRef  `json:"addedBy"`
	AddedAt   time.Time `json:"addedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SavedListItemList is a schema of the public API.
type SavedListItemList struct {
	ListID string          `json:"listId"`
	Items  []SavedListItem `json:"items"`
}

// AddSavedListItemsRequest is a schema of the public API.
type AddSavedListItemsRequest struct {
	EventIDs []int64 `json:"eventIds"`
}

// AddSavedListItemsResult is a schema of the public API.
type AddSavedListItemsResult struct {
	// Events new to the list (the others were already in it).
	Added int `json:"added"`
	// Events in the list afterwards.
	Total int `json:"total"`
}

// SavedListItemPatchClearItem is a schema of the public API.
type SavedListItemPatchClearItem string

// Values of SavedListItemPatchClearItem.
const (
	SavedListItemPatchClearItemCostBasis SavedListItemPatchClearItem = "costBasis"
	SavedListItemPatchClearItemQtyHeld   SavedListItemPatchClearItem = "qtyHeld"
)

// SavedListItemPatch is a schema of the public API.
type SavedListItemPatch struct {
	EventID int64 `json:"eventId"`
	// At most 2000 characters.
	Note *string `json:"note,omitempty"`
	// Cost per ticket in minor units.
	CostBasis *int64 `json:"costBasis,omitempty"`
	// ISO 4217 code of the cost basis.
	Currency *string `json:"currency,omitempty"`
	QtyHeld  *int    `json:"qtyHeld,omitempty"`
	// Columns to empty.
	Clear []SavedListItemPatchClearItem `json:"clear,omitempty"`
}

// UpdateSavedListItemsRequest is a schema of the public API.
type UpdateSavedListItemsRequest struct {
	Items []SavedListItemPatch `json:"items"`
}

// RemoveSavedListItemsResult is a schema of the public API.
type RemoveSavedListItemsResult struct {
	Removed int `json:"removed"`
	Total   int `json:"total"`
}

// SavedListEvent is a schema of the public API.
type SavedListEvent struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	StartsAt  *time.Time  `json:"startsAt"`
	LocalDate *string     `json:"localDate"`
	Timezone  *string     `json:"timezone"`
	Status    EventStatus `json:"status"`
	Segment   *string     `json:"segment"`
	TmEventID *string     `json:"tmEventId"`
	// The event's primary seller (`EventSummary.primarySource`).
	PrimarySource *PrimarySource  `json:"primarySource"`
	Venue         *VenueRef       `json:"venue,omitempty"`
	Attractions   []AttractionRef `json:"attractions"`
}

// OverviewHoldings: Positions recorded on list items with a quantity held, valued at the
// scope's floor (the grid's position columns, summed). Each item counts,
// so an event held in two lists counts twice. Money is per currency,
// largest first.
type OverviewHoldings struct {
	// Items with a quantity held.
	Items int `json:"items"`
	// Σ quantity held.
	Tickets int64 `json:"tickets"`
	// Σ qty × cost basis, in the cost basis's currency.
	Cost []MoneyTotal `json:"cost"`
	// Σ qty × floor, in the market's currency.
	Value []MoneyTotal `json:"value"`
	// Σ qty × (floor − cost basis) of the items with both prices in one currency.
	Pnl []MoneyTotal `json:"pnl"`
	// Σ qty × cost basis of the items in `pnl` (the base of a P&L percentage).
	PnlCost []MoneyTotal `json:"pnlCost"`
}

// Portfolio: The active org's positions (list items with a quantity held) in one market scope, per position or per event.
type Portfolio struct {
	Scope   MetricScope      `json:"scope"`
	GroupBy PortfolioGroupBy `json:"groupBy"`
	// Newest market sample among the positions.
	AsOf *time.Time `json:"asOf"`
	// The Δ24h (a fraction) at or below which a position's floor is falling (`floor_falling`).
	FloorFallingAt float64          `json:"floorFallingAt"`
	Totals         OverviewHoldings `json:"totals"`
	// One entry per bucket, in the order past, next7d, next30d, later, undated; they partition `totals`.
	Exposure []PortfolioExposure `json:"exposure"`
	// Upcoming events by start (undated last), then past events, most
	// recent first; ties by list name (by event: one row per event).
	Positions []PortfolioPosition `json:"positions"`
}

// PortfolioGroupBy: `position`: one row per list item held (an event held in two lists is
// two rows); `event`: one row per event, folding its positions across
// the org's lists.
type PortfolioGroupBy string

// Values of PortfolioGroupBy.
const (
	PortfolioGroupByPosition PortfolioGroupBy = "position"
	PortfolioGroupByEvent    PortfolioGroupBy = "event"
)

// PortfolioHistoryInterval: Bucket width of `GET /v1/portfolio/history`.
type PortfolioHistoryInterval string

// Values of PortfolioHistoryInterval.
const (
	PortfolioHistoryInterval1h PortfolioHistoryInterval = "1h"
	PortfolioHistoryInterval1d PortfolioHistoryInterval = "1d"
)

// PortfolioHistory: The portfolio's money per bucket, one series per currency (column-oriented, aligned on `t`).
type PortfolioHistory struct {
	Scope    MetricScope              `json:"scope"`
	Interval PortfolioHistoryInterval `json:"interval"`
	// The first bucket's start (after plan clamping).
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
	// Every bucket start in `[from, to)`, oldest first.
	T []time.Time `json:"t"`
	// One per currency, largest latest value first.
	Series []PortfolioHistorySeries `json:"series"`
}

// PortfolioHistorySeries: One currency's sums per bucket of `t` (minor units; null when nothing
// in the bucket has that money in this currency).
type PortfolioHistorySeries struct {
	Currency string `json:"currency"`
	// Σ qty × cost basis of the items held by the bucket's end.
	Cost []*int64 `json:"cost"`
	// Σ qty × floor.
	Value []*int64 `json:"value"`
	// Σ qty × (floor − cost basis) where floor and cost basis share this currency.
	Pnl []*int64 `json:"pnl"`
	// The cost of the items in `pnl` (P&L % = pnl ÷ pnlCost).
	PnlCost []*int64 `json:"pnlCost"`
}

// PortfolioBucket: When a position's event takes place, by its start (else the start of
// its local date, UTC): `past` once the event is past (2 hours after its
// start, the `past`/`archived` lifecycles), `next7d` within 7 days
// (started events until past), `next30d` in 7–30 days, `later` beyond,
// `undated` without a date.
type PortfolioBucket string

// Values of PortfolioBucket.
const (
	PortfolioBucketPast    PortfolioBucket = "past"
	PortfolioBucketNext7d  PortfolioBucket = "next7d"
	PortfolioBucketNext30d PortfolioBucket = "next30d"
	PortfolioBucketLater   PortfolioBucket = "later"
	PortfolioBucketUndated PortfolioBucket = "undated"
)

// PortfolioRisk: Why a position in an upcoming event is at risk: `losing_soon` — P&L
// below zero with the event within 7 days; `floor_falling` — the floor
// changed by `floorFallingAt` or less over 24h.
type PortfolioRisk string

// Values of PortfolioRisk.
const (
	PortfolioRiskLosingSoon   PortfolioRisk = "losing_soon"
	PortfolioRiskFloorFalling PortfolioRisk = "floor_falling"
)

// PortfolioExposure is a schema of the public API.
type PortfolioExposure struct {
	Bucket   PortfolioBucket  `json:"bucket"`
	Holdings OverviewHoldings `json:"holdings"`
}

// PortfolioPosition: One list item with a quantity held, valued at the scope's floor; grouped
// by event, an event's items folded (`list` null, `lists` every list
// holding it).
type PortfolioPosition struct {
	Event SavedListEvent `json:"event"`
	// The item's list; null when grouped by event.
	List *SavedListRef `json:"list"`
	// The position's list, or grouped by event every list holding the event, by name.
	Lists []SavedListRef `json:"lists"`
	// Positions of the event across the org's lists, this one included (2+ when it is held in several lists).
	EventPositions int `json:"eventPositions"`
	// Tickets held (grouped by event, summed).
	QtyHeld int `json:"qtyHeld"`
	// Cost per ticket, minor units of `costCurrency`. Grouped by event:
	// the average over the items with one (rounded); null when their
	// cost currencies differ.
	CostBasis *int64 `json:"costBasis"`
	// Currency of `costBasis` and `costTotal` (the item's, else the market's, else the face price's).
	CostCurrency *string `json:"costCurrency"`
	// qty × cost basis (grouped by event, summed over the items with a cost basis).
	CostTotal *int64 `json:"costTotal"`
	// Currency of `floor`, `marketValue` and `pnl` (the market's, else the cost basis's, else the face price's).
	Currency *string `json:"currency"`
	Floor    *int64  `json:"floor"`
	// qty × floor.
	MarketValue *int64 `json:"marketValue"`
	// qty × (floor − cost basis); null without both prices in one currency (grouped by event, summed over the items with one).
	Pnl *int64 `json:"pnl"`
	// (floor − cost basis) ÷ cost basis (grouped by event, Σ pnl ÷ the cost of the items in it).
	PnlPct *float64 `json:"pnlPct"`
	Chg24h *float64 `json:"chg24h"`
	// Tickets sold in the last 24h; null without market data.
	Sales24h *int `json:"sales24h"`
	// When the market row was computed (freshness); null without market data.
	AsOf   *time.Time      `json:"asOf"`
	Bucket PortfolioBucket `json:"bucket"`
	// At-risk reasons; empty for past events.
	Risks []PortfolioRisk `json:"risks"`
}

// EventSavedLists is a schema of the public API.
type EventSavedLists struct {
	EventID int64          `json:"eventId"`
	Items   []SavedListRef `json:"items"`
}

// PlanID: `free` is the state of an org without an active plan.
type PlanID string

// Values of PlanID.
const (
	PlanIDFree       PlanID = "free"
	PlanIDCore       PlanID = "core"
	PlanIDPro        PlanID = "pro"
	PlanIDElite      PlanID = "elite"
	PlanIDEnterprise PlanID = "enterprise"
)

// PlanRequiredProblem: A 402 problem naming the plan that allows the request.
type PlanRequiredProblem struct {
	Problem
	RequiredPlan PlanID `json:"requiredPlan"`
	CurrentPlan  PlanID `json:"currentPlan"`
	// The gated capability, e.g. `market_data`, `percent_change_filters`, `alert_rules`.
	Feature string `json:"feature"`
	// The exhausted count limit of the current plan (count features only).
	Limit *int `json:"limit,omitempty"`
}

// AlertConditionType is a schema of the public API.
type AlertConditionType string

// Values of AlertConditionType.
const (
	AlertConditionTypeMetric    AlertConditionType = "metric"
	AlertConditionTypeInventory AlertConditionType = "inventory"
	AlertConditionTypeCatalog   AlertConditionType = "catalog"
	AlertConditionTypeActivity  AlertConditionType = "activity"
	AlertConditionTypeSection   AlertConditionType = "section"
	AlertConditionTypePick      AlertConditionType = "pick"
)

// AlertActivityKind: Market activity an `activity` condition reports, one item per change
// group (section, row, price). `primary.restock`: primary seats
// restocked (released or returned, not cart holds coming back).
// `primary.sale`: an inferred primary sale (seats still gone after the
// 20-minute hold window). `primary.pulled` (opt in; needs
// `minQuantity`): Ticketmaster primary seats gone but still inside the
// hold window, about 20 minutes before they settle as a sale — they may
// be cart holds that come back. DICE ticket types (one item per type;
// `section` is the type's name, no seat count, so `minQuantity` does
// not apply): `primary.type_on_sale` a type going (back) on sale,
// `primary.type_added` a new type, `primary.waitlist` a waiting list
// opening, `primary.price_change` a new all-in price or tier
// (`listMinor`: the previous all-in price). `resale.listed`: listings
// added on Ticketmaster resale (item market `tm_resale`) or a secondary
// marketplace (item market `secondary:<marketplace>`). `resale.sold`: a
// sale on Ticketmaster resale, inferred from its listings, or on a
// secondary marketplace, confirmed by the marketplace or inferred from
// its listings.
type AlertActivityKind string

// Values of AlertActivityKind.
const (
	AlertActivityKindPrimaryRestock     AlertActivityKind = "primary.restock"
	AlertActivityKindPrimarySale        AlertActivityKind = "primary.sale"
	AlertActivityKindPrimaryPulled      AlertActivityKind = "primary.pulled"
	AlertActivityKindPrimaryTypeOnSale  AlertActivityKind = "primary.type_on_sale"
	AlertActivityKindPrimaryTypeAdded   AlertActivityKind = "primary.type_added"
	AlertActivityKindPrimaryWaitlist    AlertActivityKind = "primary.waitlist"
	AlertActivityKindPrimaryPriceChange AlertActivityKind = "primary.price_change"
	AlertActivityKindResaleListed       AlertActivityKind = "resale.listed"
	AlertActivityKindResaleSold         AlertActivityKind = "resale.sold"
)

// AlertMetric: event_metrics value a metric condition compares (EventMetrics names).
// Money metrics are in minor units; `chg*`, `depletion24h`,
// `sellThrough`, `churn24h` and `spread` are fractions (−0.15 = −15 %);
// `selloutHours` is hours (primary market only); `spread` (the floor
// over the primary floor, 0.6 = 60 % above) exists on `tm_resale` and
// secondary markets only; the rest are counts or rates.
// `chg24h`/`chg72h`/`chg7d` need the plan's % change feature.
// A `section` condition compares a section-table value instead
// (SectionRow): `floor`, `median`, `chg24h` (floorChg24h), `tickets`,
// `sales24h`, `sales72h` or `saleAvg24h`; `tickets` is section-only.
type AlertMetric string

// Values of AlertMetric.
const (
	AlertMetricFloor        AlertMetric = "floor"
	AlertMetricMedian       AlertMetric = "median"
	AlertMetricAvg24h       AlertMetric = "avg24h"
	AlertMetricSaleAvg24h   AlertMetric = "saleAvg24h"
	AlertMetricVolume24h    AlertMetric = "volume24h"
	AlertMetricChg24h       AlertMetric = "chg24h"
	AlertMetricChg72h       AlertMetric = "chg72h"
	AlertMetricChg7d        AlertMetric = "chg7d"
	AlertMetricSales24h     AlertMetric = "sales24h"
	AlertMetricSales72h     AlertMetric = "sales72h"
	AlertMetricInventory    AlertMetric = "inventory"
	AlertMetricVelocity     AlertMetric = "velocity"
	AlertMetricDepletion24h AlertMetric = "depletion24h"
	AlertMetricSellThrough  AlertMetric = "sellThrough"
	AlertMetricMovement     AlertMetric = "movement"
	AlertMetricP25          AlertMetric = "p25"
	AlertMetricSaleAvg7d    AlertMetric = "saleAvg7d"
	AlertMetricSales7d      AlertMetric = "sales7d"
	AlertMetricVolume72h    AlertMetric = "volume72h"
	AlertMetricVolume7d     AlertMetric = "volume7d"
	AlertMetricSelloutHours AlertMetric = "selloutHours"
	AlertMetricChurn24h     AlertMetric = "churn24h"
	AlertMetricSpread       AlertMetric = "spread"
	AlertMetricTickets      AlertMetric = "tickets"
)

// AlertConditionMarket: Market scope of a metric or section condition (default `primary`); secondary scopes need the secondary-marketplaces feature.
type AlertConditionMarket string

// Values of AlertConditionMarket.
const (
	AlertConditionMarketPrimary             AlertConditionMarket = "primary"
	AlertConditionMarketTmResale            AlertConditionMarket = "tm_resale"
	AlertConditionMarketSecondaryBest       AlertConditionMarket = "secondary:best"
	AlertConditionMarketSecondaryStubhub    AlertConditionMarket = "secondary:stubhub"
	AlertConditionMarketSecondaryCrowdvolt  AlertConditionMarket = "secondary:crowdvolt"
	AlertConditionMarketSecondaryTickpick   AlertConditionMarket = "secondary:tickpick"
	AlertConditionMarketSecondaryVividseats AlertConditionMarket = "secondary:vividseats"
	AlertConditionMarketSecondaryTicombo    AlertConditionMarket = "secondary:ticombo"
)

// AlertConditionOp is a schema of the public API.
type AlertConditionOp string

// Values of AlertConditionOp.
const (
	AlertConditionOpLte AlertConditionOp = "lte"
	AlertConditionOpGte AlertConditionOp = "gte"
)

// AlertConditionChange is a schema of the public API.
type AlertConditionChange string

// Values of AlertConditionChange.
const (
	AlertConditionChangeRestock   AlertConditionChange = "restock"
	AlertConditionChangeSaleBurst AlertConditionChange = "sale_burst"
)

// AlertConditionCatalog is a schema of the public API.
type AlertConditionCatalog string

// Values of AlertConditionCatalog.
const (
	AlertConditionCatalogNewEvent       AlertConditionCatalog = "new_event"
	AlertConditionCatalogNewOnsale      AlertConditionCatalog = "new_onsale"
	AlertConditionCatalogOnsaleReminder AlertConditionCatalog = "onsale_reminder"
)

// AlertConditionLeadMinutes: How long before the start an `onsale_reminder` fires (required for it).
type AlertConditionLeadMinutes int

// Values of AlertConditionLeadMinutes.
const (
	AlertConditionLeadMinutes15   AlertConditionLeadMinutes = 15
	AlertConditionLeadMinutes60   AlertConditionLeadMinutes = 60
	AlertConditionLeadMinutes180  AlertConditionLeadMinutes = 180
	AlertConditionLeadMinutes1440 AlertConditionLeadMinutes = 1440
)

// AlertConditionKind: The pick a `pick` condition reacts to.
type AlertConditionKind string

// Values of AlertConditionKind.
const (
	AlertConditionKindMustcop    AlertConditionKind = "mustcop"
	AlertConditionKindPrediction AlertConditionKind = "prediction"
)

// AlertConditionCategory: Only prediction calls of this category (`pick` with kind `prediction`; absent, both).
type AlertConditionCategory string

// Values of AlertConditionCategory.
const (
	AlertConditionCategoryTm   AlertConditionCategory = "tm"
	AlertConditionCategoryDice AlertConditionCategory = "dice"
)

// AlertConditionConfidence: Which inferred sales count (activity with `primary.sale` or `resale.sold`; default `any`, every sale).
type AlertConditionConfidence string

// Values of AlertConditionConfidence.
const (
	AlertConditionConfidenceHigh AlertConditionConfidence = "high"
	AlertConditionConfidenceAny  AlertConditionConfidence = "any"
)

// AlertCondition: What fires the rule. `metric` (event, list, attraction, venue or
// catalog filter scope): `market` `metric` `op` `value`, evaluated on
// every metrics recompute of an event in scope (catalog filters: every
// 5 min); it fires when the comparison becomes true, at most once per
// cooldown per event. `inventory` (event, list, attraction or venue
// scope): `change` `restock` (≥ `minSeats` primary seats restocked in
// one poll) or `sale_burst` (≥ `minSeats` primary seats sold within
// `windowMinutes`), optionally only in `sections` and at ≤ `maxPrice`
// all-in. `catalog` (attraction, venue or catalog filter scope):
// `catalog` `new_event` or `new_onsale` (a new event whose public
// on-sale is still ahead), or `onsale_reminder` (also list scope):
// fires `leadMinutes` before each public on-sale start (with
// `includePresales`, presale start) of a listed event in scope, once per
// event and start (checked every minute; starts already past never
// fire). `pick` (attraction, venue or catalog filter scope; the plan's
// predictions feature): fires once per event when a `kind` `mustcop` or
// `prediction` call opens on it (predictions: optionally only of
// `category` and at most `maxRisk`), provided the org sees those picks
// (published to it); the message carries the not-financial-advice
// notice. `activity` (any scope): every change of the
// `activities` kinds that passes the filters — `markets`, `sections`,
// `minQuantity` seats per change, all-in `minPrice`/`maxPrice` and, for
// `primary.sale` and marketplace `resale.sold`, `confidence` (`any`:
// every sale, the default; `high`: only primary sales seen between
// close polls and marketplace sales confirmed or inferred with high
// confidence) — is
// reported, no cooldown; bursts are merged by the rule's
// `digestSeconds`. A catalog filter scope needs at least one of these
// filters and a digest window. `section` (event, list, attraction or
// venue scope): `market` `metric` `op` `value` over each section of
// the event's section table (`GET /v1/events/{id}/sections`),
// optionally only for `sections` (section keys or labels, matched
// case-insensitively; "Section 112" matches `112`), evaluated on every
// metrics recompute of an event in scope; it fires when the comparison
// becomes true for a section, at most once per cooldown per event and
// section (several sections at once: one digest). Price metrics
// (`floor`, `median`, `chg24h`, `saleAvg24h`) need the plan's map
// prices, `chg24h` also its % change feature, secondary markets the
// secondary-marketplaces feature.
type AlertCondition struct {
	Type AlertConditionType `json:"type"`
	// Market scope of a metric or section condition (default `primary`); secondary scopes need the secondary-marketplaces feature.
	Market   *AlertConditionMarket `json:"market,omitempty"`
	Metric   *AlertMetric          `json:"metric,omitempty"`
	Op       *AlertConditionOp     `json:"op,omitempty"`
	Value    *float64              `json:"value,omitempty"`
	Change   *AlertConditionChange `json:"change,omitempty"`
	MinSeats *int                  `json:"minSeats,omitempty"`
	// Section names (at most 50, each at most 100 characters; more is 400).
	Sections []string `json:"sections,omitempty"`
	// All-in price ceiling in minor units.
	MaxPrice      *int64                 `json:"maxPrice,omitempty"`
	WindowMinutes *int                   `json:"windowMinutes,omitempty"`
	Catalog       *AlertConditionCatalog `json:"catalog,omitempty"`
	// How long before the start an `onsale_reminder` fires (required for it).
	LeadMinutes *AlertConditionLeadMinutes `json:"leadMinutes,omitempty"`
	// An `onsale_reminder` also fires for presale starts.
	IncludePresales *bool `json:"includePresales,omitempty"`
	// The pick a `pick` condition reacts to.
	Kind *AlertConditionKind `json:"kind,omitempty"`
	// Only prediction calls of this category (`pick` with kind `prediction`; absent, both).
	Category *AlertConditionCategory `json:"category,omitempty"`
	// Only prediction calls with at most this risk (`pick` with kind `prediction`).
	MaxRisk    *int                `json:"maxRisk,omitempty"`
	Activities []AlertActivityKind `json:"activities,omitempty"`
	// Only changes of at least this many seats (activity; default 1, required with `primary.pulled`; not applied to DICE ticket types).
	MinQuantity *int `json:"minQuantity,omitempty"`
	// All-in price floor in minor units (activity).
	MinPrice *int64 `json:"minPrice,omitempty"`
	// Which inferred sales count (activity with `primary.sale` or `resale.sold`; default `any`, every sale).
	Confidence *AlertConditionConfidence `json:"confidence,omitempty"`
	// Only changes in these markets (activity; empty or absent: every
	// market). Restocks, primary sales, pulled seats and DICE ticket types are primary-market only.
	// Secondary markets need the secondary-marketplaces feature; without
	// it, a rule without `markets` sees Ticketmaster changes only.
	Markets []WebhookMarket `json:"markets,omitempty"`
}

// AlertCatalogFilter: Upcoming events matching every given field (keyword matches the folded event name).
type AlertCatalogFilter struct {
	// ISO 3166-1 alpha-2.
	Country *string `json:"country,omitempty"`
	City    *string `json:"city,omitempty"`
	Segment *string `json:"segment,omitempty"`
	Genre   *string `json:"genre,omitempty"`
	VenueID *int64  `json:"venueId,omitempty"`
	Keyword *string `json:"keyword,omitempty"`
}

// AlertScopeKind is a schema of the public API.
type AlertScopeKind string

// Values of AlertScopeKind.
const (
	AlertScopeKindEvent         AlertScopeKind = "event"
	AlertScopeKindList          AlertScopeKind = "list"
	AlertScopeKindAttraction    AlertScopeKind = "attraction"
	AlertScopeKindVenue         AlertScopeKind = "venue"
	AlertScopeKindCatalogFilter AlertScopeKind = "catalog_filter"
)

// AlertScopeInput is a schema of the public API.
type AlertScopeInput struct {
	Kind AlertScopeKind `json:"kind"`
	// Event id, list uuid, attraction id or venue id (not for catalog_filter).
	Ref    *string             `json:"ref,omitempty"`
	Filter *AlertCatalogFilter `json:"filter,omitempty"`
}

// AlertScope is a schema of the public API.
type AlertScope struct {
	Kind   AlertScopeKind      `json:"kind"`
	Ref    *string             `json:"ref,omitempty"`
	Filter *AlertCatalogFilter `json:"filter,omitempty"`
	// Display name (event, list or attraction name, or the filter summary).
	Label string `json:"label"`
}

// AlertRuleInput is a schema of the public API.
type AlertRuleInput struct {
	Name      string          `json:"name"`
	Scope     AlertScopeInput `json:"scope"`
	Condition AlertCondition  `json:"condition"`
	// Metric and inventory rules (activity and catalog rules ignore it).
	CooldownMinutes *int                `json:"cooldownMinutes,omitempty"`
	DigestSeconds   *AlertDigestSeconds `json:"digestSeconds,omitempty"`
	// Notification channel ids of the org (at most the org's channel cap plus 25, 50 by default; more is 400).
	Channels []int64 `json:"channels,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

// AlertRulePatch is a schema of the public API.
type AlertRulePatch struct {
	Name            *string             `json:"name,omitempty"`
	Scope           *AlertScopeInput    `json:"scope,omitempty"`
	Condition       *AlertCondition     `json:"condition,omitempty"`
	CooldownMinutes *int                `json:"cooldownMinutes,omitempty"`
	DigestSeconds   *AlertDigestSeconds `json:"digestSeconds,omitempty"`
	// Notification channel ids of the org (at most the org's channel cap plus 25, 50 by default; more is 400).
	Channels []int64 `json:"channels,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

// AlertDigestSeconds: Digest window: firings of the rule within this many seconds are sent
// to each channel as one message (0: at once). Defaults to 30 for
// activity rules and 0 otherwise; catalog-filter activity rules need at
// least 30. Whatever the window, firings waiting on a rate-limited
// channel are merged into its next message.
type AlertDigestSeconds int

// Values of AlertDigestSeconds.
const (
	AlertDigestSeconds0   AlertDigestSeconds = 0
	AlertDigestSeconds30  AlertDigestSeconds = 30
	AlertDigestSeconds60  AlertDigestSeconds = 60
	AlertDigestSeconds300 AlertDigestSeconds = 300
)

// AlertRule is a schema of the public API.
type AlertRule struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Scope     AlertScope     `json:"scope"`
	Condition AlertCondition `json:"condition"`
	// One-line summary of the condition ("Primary floor ≤ $120").
	Description     string             `json:"description"`
	CooldownMinutes int                `json:"cooldownMinutes"`
	DigestSeconds   AlertDigestSeconds `json:"digestSeconds"`
	Channels        []int64            `json:"channels"`
	Enabled         bool               `json:"enabled"`
	// The rule sends nothing until then (`PUT /v1/alerts/rules/{ruleId}/snooze`); absent when not snoozed or the snooze has ended.
	SnoozedUntil *time.Time     `json:"snoozedUntil,omitempty"`
	LastFiredAt  *time.Time     `json:"lastFiredAt,omitempty"`
	FireCount    int64          `json:"fireCount"`
	LastSkip     *AlertRuleSkip `json:"lastSkip,omitempty"`
	CreatedBy    *string        `json:"createdBy,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

// AlertRuleSkipReason is a schema of the public API.
type AlertRuleSkipReason string

// Values of AlertRuleSkipReason.
const (
	AlertRuleSkipReasonPlan     AlertRuleSkipReason = "plan"
	AlertRuleSkipReasonCooldown AlertRuleSkipReason = "cooldown"
	AlertRuleSkipReasonSnoozed  AlertRuleSkipReason = "snoozed"
)

// AlertRuleSkip: Why the rule last matched without firing; absent once it has fired
// since, and after its scope or condition changes. `snoozed`: the rule
// is snoozed. `plan`: the org's current plan doesn't allow the rule (a
// feature it lacks, or the rule is beyond the plan's rule count;
// dropped as soon as the plan allows it again). `cooldown`: it already
// fired for that event less than its cooldown ago. `at` is refreshed at
// most every 5 minutes while the reason stays the same.
type AlertRuleSkip struct {
	// When the rule matched without firing.
	At     time.Time           `json:"at"`
	Reason AlertRuleSkipReason `json:"reason"`
	// Human explanation ("Your plan doesn't include market data (needs Core).").
	Detail string `json:"detail"`
	// `cooldown`, `snoozed`: when the cooldown or snooze ends.
	Until *time.Time `json:"until,omitempty"`
	// `plan`: the gated capability, as in PlanRequiredProblem (`market_data`, `alert_rules`, …).
	Feature      *string `json:"feature,omitempty"`
	RequiredPlan *PlanID `json:"requiredPlan,omitempty"`
}

// AlertRuleSnoozeInput is a schema of the public API.
type AlertRuleSnoozeInput struct {
	// When the snooze ends (in the future, at most 90 days ahead).
	Until time.Time `json:"until"`
}

// AlertRuleUsage is a schema of the public API.
type AlertRuleUsage struct {
	// Rules the org has.
	Rules int `json:"rules"`
	// The plan's limit; absent when unlimited.
	MaxRules *int `json:"maxRules,omitempty"`
}

// AlertRulePage is a schema of the public API.
type AlertRulePage struct {
	Items      []AlertRule    `json:"items"`
	NextBefore *int64         `json:"nextBefore,omitempty"`
	Usage      AlertRuleUsage `json:"usage"`
}

// AlertField is a schema of the public API.
type AlertField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AlertDeliveryStatus is a schema of the public API.
type AlertDeliveryStatus string

// Values of AlertDeliveryStatus.
const (
	AlertDeliveryStatusPending  AlertDeliveryStatus = "pending"
	AlertDeliveryStatusSent     AlertDeliveryStatus = "sent"
	AlertDeliveryStatusRetrying AlertDeliveryStatus = "retrying"
	AlertDeliveryStatusFailed   AlertDeliveryStatus = "failed"
	AlertDeliveryStatusSkipped  AlertDeliveryStatus = "skipped"
)

// AlertDelivery is a schema of the public API.
type AlertDelivery struct {
	ID int64 `json:"id"`
	// Null once the channel was deleted.
	ChannelID   *int64              `json:"channelId"`
	ChannelKind ChannelKind         `json:"channelKind"`
	ChannelName string              `json:"channelName"`
	Status      AlertDeliveryStatus `json:"status"`
	Attempts    int                 `json:"attempts"`
	Error       *string             `json:"error"`
	SentAt      *time.Time          `json:"sentAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	// When the delivery is sent at the earliest (the end of the rule's digest window).
	DueAt time.Time `json:"dueAt"`
	// The delivery whose message carried this one (merged by the digest window or while the channel was rate-limited).
	CoalescedInto *int64 `json:"coalescedInto"`
}

// AlertEventTone is a schema of the public API.
type AlertEventTone string

// Values of AlertEventTone.
const (
	AlertEventToneUp   AlertEventTone = "up"
	AlertEventToneDown AlertEventTone = "down"
)

// AlertEvent: One firing. Also the `data` of the `alerts` stream topic. Digest
// firings (several events newly matched at once) have no `event` and
// list them in `events`. Activity firings list what happened in
// `items` (at most 50; `itemsTotal` counts them all).
type AlertEvent struct {
	ID int64 `json:"id"`
	// Null once the rule was deleted.
	RuleID   *int64             `json:"ruleId"`
	RuleName string             `json:"ruleName"`
	Kind     AlertConditionType `json:"kind"`
	Test     bool               `json:"test"`
	FiredAt  time.Time          `json:"firedAt"`
	Title    string             `json:"title"`
	Body     string             `json:"body"`
	// Link into the app (event page, or the rule's history for digests).
	URL string `json:"url"`
	// The event's page at its primary seller (activity firings).
	TicketURL *string `json:"ticketUrl,omitempty"`
	// The seller of `ticketUrl`; absent when it is neither Ticketmaster nor a primary source.
	TicketSource *PrimarySource      `json:"ticketSource,omitempty"`
	Tone         *AlertEventTone     `json:"tone,omitempty"`
	Event        *SaleEventRef       `json:"event,omitempty"`
	Events       []SaleEventRef      `json:"events,omitempty"`
	Fields       []AlertField        `json:"fields"`
	Items        []AlertActivityItem `json:"items,omitempty"`
	// Changes the firing covered (activity; `items` lists the first 50).
	ItemsTotal *int `json:"itemsTotal,omitempty"`
	// Raw values behind the firing (metric, value, threshold, seats, …).
	Data       map[string]any  `json:"data,omitempty"`
	Deliveries []AlertDelivery `json:"deliveries"`
}

// AlertActivityItem: One change an activity firing reports (a section, row and price group of seats).
type AlertActivityItem struct {
	Activity AlertActivityKind `json:"activity"`
	// `primary`, `tm_resale` or `secondary:<marketplace>`.
	Market  string  `json:"market"`
	Section string  `json:"section"`
	Row     *string `json:"row,omitempty"`
	// Seat labels (absent for general admission).
	Seats    []string `json:"seats,omitempty"`
	Quantity int      `json:"quantity"`
	// All-in price per ticket in minor units (absent when unpriced).
	PriceMinor *int64 `json:"priceMinor,omitempty"`
	// Face (list) price per ticket in minor units.
	ListMinor *int64  `json:"listMinor,omitempty"`
	Currency  *string `json:"currency,omitempty"`
	// For sales, `inferred_high` or `inferred_low` (`resale.sold` is always `inferred_low`).
	Confidence *string   `json:"confidence,omitempty"`
	At         time.Time `json:"at"`
}

// AlertEventPage is a schema of the public API.
type AlertEventPage struct {
	Items      []AlertEvent `json:"items"`
	NextBefore *int64       `json:"nextBefore,omitempty"`
}

// AlertPreviewItem is a schema of the public API.
type AlertPreviewItem struct {
	Event SaleEventRef `json:"event"`
	// The current value, or what happened (restock size, announcement time).
	Detail string     `json:"detail"`
	At     *time.Time `json:"at,omitempty"`
}

// AlertPreview is a schema of the public API.
type AlertPreview struct {
	Items []AlertPreviewItem `json:"items"`
	Total int                `json:"total"`
	// What the items cover ("matching now", "restocks in the last 24 hours", …).
	Window string `json:"window"`
}

// ChannelKind: `slack_app` channels post as the org's installed Slack app (Settings →
// Integrations; created with `POST /v1/integrations/slack/workspaces/{teamId}/channels`,
// not `POST /v1/channels`). `webpush` and `apns` channels are members'
// own browsers and iOS devices.
type ChannelKind string

// Values of ChannelKind.
const (
	ChannelKindDiscordWebhook    ChannelKind = "discord_webhook"
	ChannelKindSlackWebhook      ChannelKind = "slack_webhook"
	ChannelKindSlackApp          ChannelKind = "slack_app"
	ChannelKindTeamsWebhook      ChannelKind = "teams_webhook"
	ChannelKindGoogleChatWebhook ChannelKind = "google_chat_webhook"
	ChannelKindTelegram          ChannelKind = "telegram"
	ChannelKindEmail             ChannelKind = "email"
	ChannelKindWebpush           ChannelKind = "webpush"
	ChannelKindApns              ChannelKind = "apns"
)

// NotificationChannel is a schema of the public API.
type NotificationChannel struct {
	ID   int64       `json:"id"`
	Kind ChannelKind `json:"kind"`
	Name string      `json:"name"`
	// Non-secret destination hint (webhook host, chat title, address, push service, APNs environment).
	Label    string `json:"label"`
	Verified bool   `json:"verified"`
	Enabled  bool   `json:"enabled"`
	// An email code or Telegram link is outstanding.
	VerificationPending bool       `json:"verificationPending"`
	VerifyExpiresAt     *time.Time `json:"verifyExpiresAt,omitempty"`
	// A device channel (Web Push browser, iOS device) of the signed-in caller; always false for API keys.
	Mine        bool       `json:"mine"`
	LastError   *string    `json:"lastError,omitempty"`
	LastErrorAt *time.Time `json:"lastErrorAt,omitempty"`
	LastSentAt  *time.Time `json:"lastSentAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// NotificationChannelList is a schema of the public API.
type NotificationChannelList struct {
	Items []NotificationChannel `json:"items"`
}

// NotificationChannelCreated is a schema of the public API.
type NotificationChannelCreated struct {
	Channel NotificationChannel `json:"channel"`
	// Telegram deep link (t.me/<bot>?start=<code>) that connects the chat it is opened in; shown once, valid 1 h.
	LinkURL *string `json:"linkUrl,omitempty"`
}

// CreateChannelRequestKind is a schema of the public API.
type CreateChannelRequestKind string

// Values of CreateChannelRequestKind.
const (
	CreateChannelRequestKindDiscordWebhook    CreateChannelRequestKind = "discord_webhook"
	CreateChannelRequestKindSlackWebhook      CreateChannelRequestKind = "slack_webhook"
	CreateChannelRequestKindTeamsWebhook      CreateChannelRequestKind = "teams_webhook"
	CreateChannelRequestKindGoogleChatWebhook CreateChannelRequestKind = "google_chat_webhook"
	CreateChannelRequestKindTelegram          CreateChannelRequestKind = "telegram"
	CreateChannelRequestKindEmail             CreateChannelRequestKind = "email"
)

// CreateChannelRequest is a schema of the public API.
type CreateChannelRequest struct {
	Kind CreateChannelRequestKind `json:"kind"`
	Name string                   `json:"name"`
	// Incoming-webhook URL (discord_webhook, slack_webhook, teams_webhook, google_chat_webhook).
	URL *string `json:"url,omitempty"`
	// Email address (email).
	Address *string `json:"address,omitempty"`
}

// UpdateChannelRequest is a schema of the public API.
type UpdateChannelRequest struct {
	Name    *string `json:"name,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
}

// VerifyChannelRequest is a schema of the public API.
type VerifyChannelRequest struct {
	Code string `json:"code"`
}

// UsageGroup: Endpoint group a request is metered under.
type UsageGroup string

// Values of UsageGroup.
const (
	UsageGroupCatalog  UsageGroup = "catalog"
	UsageGroupMarket   UsageGroup = "market"
	UsageGroupSales    UsageGroup = "sales"
	UsageGroupMaps     UsageGroup = "maps"
	UsageGroupLists    UsageGroup = "lists"
	UsageGroupAlerts   UsageGroup = "alerts"
	UsageGroupWebhooks UsageGroup = "webhooks"
	UsageGroupExports  UsageGroup = "exports"
	UsageGroupStream   UsageGroup = "stream"
	UsageGroupUsage    UsageGroup = "usage"
)

// UsageQuota: The plan's API limits (0 = unlimited).
type UsageQuota struct {
	RequestsPerMinute   int `json:"requestsPerMinute"`
	RequestsPerDay      int `json:"requestsPerDay"`
	ConcurrentRequests  int `json:"concurrentRequests"`
	StreamConnections   int `json:"streamConnections"`
	StreamSubscriptions int `json:"streamSubscriptions"`
}

// UsageCurrent: Live usage against the quota.
type UsageCurrent struct {
	RequestsThisMinute int `json:"requestsThisMinute"`
	// Requests since 00:00 UTC.
	RequestsToday int `json:"requestsToday"`
	// Open stream connections of the org's keys.
	StreamConnections int `json:"streamConnections"`
}

// UsageCounts is a schema of the public API.
type UsageCounts struct {
	// Requests answered (including 429s).
	Requests int64 `json:"requests"`
	// Requests refused with 429.
	Throttled int64 `json:"throttled"`
	// Requests answered with a 5xx.
	Errors int64 `json:"errors"`
	// Stream connections opened.
	StreamConnections int64 `json:"streamConnections"`
	// Stream messages delivered.
	StreamMessages int64 `json:"streamMessages"`
}

// UsageBucket is a schema of the public API.
type UsageBucket struct {
	UsageCounts
	// Bucket start.
	Ts time.Time `json:"ts"`
}

// UsageGroupCounts is a schema of the public API.
type UsageGroupCounts struct {
	UsageCounts
	Group UsageGroup `json:"group"`
}

// UsageKeyCounts is a schema of the public API.
type UsageKeyCounts struct {
	UsageCounts
	APIKeyID string `json:"apiKeyId"`
	KeyID    string `json:"keyId"`
	Name     string `json:"name"`
	// Requests per bucket, aligned with `buckets`.
	Series []int64 `json:"series"`
}

// UsageGranularity is a schema of the public API.
type UsageGranularity string

// Values of UsageGranularity.
const (
	UsageGranularityMinute UsageGranularity = "minute"
	UsageGranularityHour   UsageGranularity = "hour"
	UsageGranularityDay    UsageGranularity = "day"
)

// Usage is a schema of the public API.
type Usage struct {
	From        time.Time          `json:"from"`
	To          time.Time          `json:"to"`
	Granularity UsageGranularity   `json:"granularity"`
	Quota       UsageQuota         `json:"quota"`
	Current     UsageCurrent       `json:"current"`
	Totals      UsageCounts        `json:"totals"`
	Buckets     []UsageBucket      `json:"buckets"`
	Groups      []UsageGroupCounts `json:"groups"`
	Keys        []UsageKeyCounts   `json:"keys"`
}

// RateLimitedProblemLimit: The exhausted limit (`key_minute` / `key_day`, the API key's own quotas).
type RateLimitedProblemLimit string

// Values of RateLimitedProblemLimit.
const (
	RateLimitedProblemLimitMinute            RateLimitedProblemLimit = "minute"
	RateLimitedProblemLimitDay               RateLimitedProblemLimit = "day"
	RateLimitedProblemLimitConcurrent        RateLimitedProblemLimit = "concurrent"
	RateLimitedProblemLimitStreamConnections RateLimitedProblemLimit = "stream_connections"
	RateLimitedProblemLimitUser              RateLimitedProblemLimit = "user"
	RateLimitedProblemLimitIP                RateLimitedProblemLimit = "ip"
	RateLimitedProblemLimitStreamTickets     RateLimitedProblemLimit = "stream_tickets"
	RateLimitedProblemLimitKeyMinute         RateLimitedProblemLimit = "key_minute"
	RateLimitedProblemLimitKeyDay            RateLimitedProblemLimit = "key_day"
)

// RateLimitedProblem is a schema of the public API.
type RateLimitedProblem struct {
	Problem
	// The exhausted limit (`key_minute` / `key_day`, the API key's own quotas).
	Limit RateLimitedProblemLimit `json:"limit"`
}

// WebhookEventType: `alert.fired`: an alert rule of the org fired (payload: the firing);
// `onsale.announced`: a newly announced event, on-sale or presale (the
// `onsales` stream payload); `event.status_changed`: an event of the
// endpoint's list changed status (onsale, offsale, canceled, postponed,
// rescheduled); `list.sale`: a sale on an event of the endpoint's list
// (the `Sale` schema). Market activity on events of the endpoint's list
// (payload: `WebhookActivity`, one delivery per poll of the event):
// `inventory.restocked`: primary seats became available (restocks,
// GA restocks); `inventory.sold`: settled primary sales (inferred);
// `resale.listed`: listings were added on Ticketmaster resale;
// `resale.sold`: Ticketmaster resale sales inferred from its listings
// (a listing sold down or vanished; each change carries its
// confidence); `marketplace.listed`: listings were added on a secondary
// marketplace; `marketplace.sold`: secondary-marketplace sales,
// confirmed by the marketplace or inferred from its listings (each
// change carries its confidence).
type WebhookEventType string

// Values of WebhookEventType.
const (
	WebhookEventTypeAlertFired         WebhookEventType = "alert.fired"
	WebhookEventTypeOnsaleAnnounced    WebhookEventType = "onsale.announced"
	WebhookEventTypeEventStatusChanged WebhookEventType = "event.status_changed"
	WebhookEventTypeListSale           WebhookEventType = "list.sale"
	WebhookEventTypeInventoryRestocked WebhookEventType = "inventory.restocked"
	WebhookEventTypeInventorySold      WebhookEventType = "inventory.sold"
	WebhookEventTypeResaleListed       WebhookEventType = "resale.listed"
	WebhookEventTypeResaleSold         WebhookEventType = "resale.sold"
	WebhookEventTypeMarketplaceListed  WebhookEventType = "marketplace.listed"
	WebhookEventTypeMarketplaceSold    WebhookEventType = "marketplace.sold"
)

// WebhookMarket: A market of an endpoint's `markets` filter: `primary` (Ticketmaster
// primary), `tm_resale` (Ticketmaster resale) or
// `secondary:<marketplace>` (`list.sale` and marketplace activity of a
// marketplace; needs the secondary-marketplaces feature to deliver).
// More marketplaces may be added; treat unknown values as other markets.
type WebhookMarket string

// Values of WebhookMarket.
const (
	WebhookMarketPrimary             WebhookMarket = "primary"
	WebhookMarketTmResale            WebhookMarket = "tm_resale"
	WebhookMarketSecondaryStubhub    WebhookMarket = "secondary:stubhub"
	WebhookMarketSecondaryCrowdvolt  WebhookMarket = "secondary:crowdvolt"
	WebhookMarketSecondaryTickpick   WebhookMarket = "secondary:tickpick"
	WebhookMarketSecondaryVividseats WebhookMarket = "secondary:vividseats"
	WebhookMarketSecondaryTicombo    WebhookMarket = "secondary:ticombo"
)

// WebhookStatus is a schema of the public API.
type WebhookStatus string

// Values of WebhookStatus.
const (
	WebhookStatusActive   WebhookStatus = "active"
	WebhookStatusDisabled WebhookStatus = "disabled"
)

// Webhook is a schema of the public API.
type Webhook struct {
	ID          string             `json:"id"`
	URL         string             `json:"url"`
	Description string             `json:"description"`
	Events      []WebhookEventType `json:"events"`
	// Scopes `list.sale`, `event.status_changed` and the market-activity types to this list's events.
	ListID *string `json:"listId"`
	// Markets `list.sale` and the market-activity types are limited to; empty means every market.
	Markets []WebhookMarket `json:"markets"`
	Status  WebhookStatus   `json:"status"`
	// Why the endpoint is disabled (auto-disabled after 24 h of failed deliveries, or by a user).
	DisabledReason      *string `json:"disabledReason"`
	ConsecutiveFailures int     `json:"consecutiveFailures"`
	// First failure of the current failure streak.
	FailingSince   *time.Time `json:"failingSince"`
	LastDeliveryAt *time.Time `json:"lastDeliveryAt"`
	LastSuccessAt  *time.Time `json:"lastSuccessAt"`
	// The secret's last four characters, `whsec_…abcd`.
	SecretHint string    `json:"secretHint"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  *UserRef  `json:"createdBy"`
}

// WebhookList is a schema of the public API.
type WebhookList struct {
	Items []Webhook `json:"items"`
}

// WebhookInput is a schema of the public API.
type WebhookInput struct {
	URL         string             `json:"url"`
	Description *string            `json:"description,omitempty"`
	Events      []WebhookEventType `json:"events"`
	ListID      Nullable[string]   `json:"listId,omitzero"`
	// Limit `list.sale` and the market-activity types to these markets; omitted or empty means every market.
	Markets []WebhookMarket `json:"markets,omitempty"`
}

// WebhookPatch is a schema of the public API.
type WebhookPatch struct {
	URL         *string            `json:"url,omitempty"`
	Description *string            `json:"description,omitempty"`
	Events      []WebhookEventType `json:"events,omitempty"`
	ListID      Nullable[string]   `json:"listId,omitzero"`
	// Replaces the market filter; empty means every market.
	Markets []WebhookMarket `json:"markets,omitempty"`
	// Remove the list filter.
	ClearList *bool          `json:"clearList,omitempty"`
	Status    *WebhookStatus `json:"status,omitempty"`
}

// WebhookSecret is a schema of the public API.
type WebhookSecret struct {
	Webhook Webhook `json:"webhook"`
	// The signing secret (`whsec_` + 32 random bytes in base62). The HMAC key is this whole string as UTF-8 bytes.
	Secret string `json:"secret"`
}

// WebhookTestInput is a schema of the public API.
type WebhookTestInput struct {
	Type *WebhookEventType `json:"type,omitempty"`
}

// WebhookDeliveryStatus: `pending`: queued or retrying; `succeeded`: a 2xx answer; `failed`: gave up (24 h of retries, endpoint disabled or deleted, or a permanent error).
type WebhookDeliveryStatus string

// Values of WebhookDeliveryStatus.
const (
	WebhookDeliveryStatusPending   WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusSucceeded WebhookDeliveryStatus = "succeeded"
	WebhookDeliveryStatusFailed    WebhookDeliveryStatus = "failed"
)

// WebhookDelivery is a schema of the public API.
type WebhookDelivery struct {
	// This delivery (`Status302-Delivery-Id`).
	ID        string `json:"id"`
	WebhookID string `json:"webhookId"`
	// The envelope `id`; redeliveries keep it, so receivers can dedupe on it.
	EventID        string                `json:"eventId"`
	Event          WebhookEventType      `json:"event"`
	Test           bool                  `json:"test"`
	Status         WebhookDeliveryStatus `json:"status"`
	Attempts       int                   `json:"attempts"`
	ResponseStatus *int                  `json:"responseStatus"`
	// First 512 bytes of the last response body.
	ResponseSnippet *string `json:"responseSnippet"`
	// Duration of the last attempt.
	LatencyMs     *int       `json:"latencyMs"`
	Error         *string    `json:"error"`
	CreatedAt     time.Time  `json:"createdAt"`
	LastAttemptAt *time.Time `json:"lastAttemptAt"`
	NextAttemptAt *time.Time `json:"nextAttemptAt"`
	RedeliveryOf  *string    `json:"redeliveryOf"`
}

// WebhookDeliveryPage is a schema of the public API.
type WebhookDeliveryPage struct {
	Items      []WebhookDelivery `json:"items"`
	NextCursor *string           `json:"nextCursor"`
}

// WebhookEnvelope: The JSON body of every webhook delivery. `data` depends on `type`:
// `WebhookActivity` for the market-activity types, the `Sale` schema for
// `list.sale` (see `WebhookEventType`).
type WebhookEnvelope struct {
	// The event id, kept by retries and redeliveries; process each id once.
	ID        string           `json:"id"`
	Type      WebhookEventType `json:"type"`
	CreatedAt time.Time        `json:"createdAt"`
	Test      bool             `json:"test"`
	Data      map[string]any   `json:"data"`
}

// WebhookActivityEnvelopeType is a schema of the public API.
type WebhookActivityEnvelopeType string

// Values of WebhookActivityEnvelopeType.
const (
	WebhookActivityEnvelopeTypeInventoryRestocked WebhookActivityEnvelopeType = "inventory.restocked"
	WebhookActivityEnvelopeTypeInventorySold      WebhookActivityEnvelopeType = "inventory.sold"
	WebhookActivityEnvelopeTypeResaleListed       WebhookActivityEnvelopeType = "resale.listed"
	WebhookActivityEnvelopeTypeResaleSold         WebhookActivityEnvelopeType = "resale.sold"
	WebhookActivityEnvelopeTypeMarketplaceListed  WebhookActivityEnvelopeType = "marketplace.listed"
	WebhookActivityEnvelopeTypeMarketplaceSold    WebhookActivityEnvelopeType = "marketplace.sold"
)

// WebhookActivityEnvelope: The envelope of `inventory.restocked`, `inventory.sold`, `resale.listed`, `resale.sold`, `marketplace.listed` and `marketplace.sold` deliveries.
type WebhookActivityEnvelope struct {
	// Derived from the source poll and the type, so a poll is delivered
	// at most once per type even across retries and restarts.
	ID        string                      `json:"id"`
	Type      WebhookActivityEnvelopeType `json:"type"`
	CreatedAt time.Time                   `json:"createdAt"`
	Test      bool                        `json:"test"`
	Data      WebhookActivity             `json:"data"`
}

// WebhookActivity: The changes of one type found by one poll of an event. Changes of a
// mass transition (a hold pulled or an event withdrawn in one poll) are
// never delivered.
type WebhookActivity struct {
	Event  WebhookEventRef `json:"event"`
	Market WebhookMarket   `json:"market"`
	// Time of the poll that saw the changes.
	At time.Time `json:"at"`
	// Tickets over every change of the poll, including any left out of `changes`.
	Quantity int `json:"quantity"`
	// More than 200 changes; `changes` holds the first 200.
	Truncated bool                    `json:"truncated"`
	Changes   []WebhookActivityChange `json:"changes"`
}

// WebhookActivityChangeKind: `restock`/`ga_restock` (inventory.restocked), `sale`/`ga_sale`
// (inventory.sold), `added` (resale.listed, marketplace.listed),
// `confirmed` (reported by the marketplace), `partial` (a
// listing's quantity dropped) or `removal` (a listing vanished and
// stayed gone) (resale.sold, marketplace.sold).
type WebhookActivityChangeKind string

// Values of WebhookActivityChangeKind.
const (
	WebhookActivityChangeKindRestock   WebhookActivityChangeKind = "restock"
	WebhookActivityChangeKindGaRestock WebhookActivityChangeKind = "ga_restock"
	WebhookActivityChangeKindSale      WebhookActivityChangeKind = "sale"
	WebhookActivityChangeKindGaSale    WebhookActivityChangeKind = "ga_sale"
	WebhookActivityChangeKindAdded     WebhookActivityChangeKind = "added"
	WebhookActivityChangeKindConfirmed WebhookActivityChangeKind = "confirmed"
	WebhookActivityChangeKindPartial   WebhookActivityChangeKind = "partial"
	WebhookActivityChangeKindRemoval   WebhookActivityChangeKind = "removal"
)

// WebhookActivityChange: Seats of one section, row and price that changed together; GA pools
// have an empty row and no seats. Marketplace listings carry no seat
// numbers (empty `seats`).
type WebhookActivityChange struct {
	// `restock`/`ga_restock` (inventory.restocked), `sale`/`ga_sale`
	// (inventory.sold), `added` (resale.listed, marketplace.listed),
	// `confirmed` (reported by the marketplace), `partial` (a
	// listing's quantity dropped) or `removal` (a listing vanished and
	// stayed gone) (resale.sold, marketplace.sold).
	Kind     WebhookActivityChangeKind `json:"kind"`
	Section  string                    `json:"section"`
	Row      string                    `json:"row"`
	Seats    []string                  `json:"seats"`
	Quantity int                       `json:"quantity"`
	// All-in price per ticket in minor units; null when unpriced.
	PriceMinor *int64 `json:"priceMinor"`
	// Face (list) price per ticket before fees.
	ListMinor *int64  `json:"listMinor"`
	Currency  *string `json:"currency"`
	// Sales only (null otherwise): primary sales are `inferred_high`
	// when the polls around them were close together, else
	// `inferred_low`; resale sales are always `inferred_low`;
	// marketplace sales are `confirmed`, or inferred like primary
	// sales.
	Confidence *SaleConfidence `json:"confidence"`
}

// WebhookEventRef: The event summary of a market-activity payload.
type WebhookEventRef struct {
	SaleEventRef
	// The event's Ticketmaster page, when known.
	URL *string `json:"url"`
}

// ExportKind is a schema of the public API.
type ExportKind string

// Values of ExportKind.
const (
	ExportKindEvents ExportKind = "events"
	ExportKindSales  ExportKind = "sales"
	ExportKindSeries ExportKind = "series"
	ExportKindList   ExportKind = "list"
)

// ExportFormat is a schema of the public API.
type ExportFormat string

// Values of ExportFormat.
const (
	ExportFormatCSV   ExportFormat = "csv"
	ExportFormatJsonl ExportFormat = "jsonl"
)

// ExportStatus is a schema of the public API.
type ExportStatus string

// Values of ExportStatus.
const (
	ExportStatusQueued    ExportStatus = "queued"
	ExportStatusRunning   ExportStatus = "running"
	ExportStatusSucceeded ExportStatus = "succeeded"
	ExportStatusFailed    ExportStatus = "failed"
	ExportStatusExpired   ExportStatus = "expired"
)

// ExportInput is a schema of the public API.
type ExportInput struct {
	Kind   ExportKind   `json:"kind"`
	Format ExportFormat `json:"format"`
	// Query parameters of the mirrored read (strings, numbers, booleans or arrays of them), plus `eventId`/`listId` where the read has them in its path.
	Filters map[string]any `json:"filters,omitempty"`
}

// Export is a schema of the public API.
type Export struct {
	ID      string         `json:"id"`
	Kind    ExportKind     `json:"kind"`
	Format  ExportFormat   `json:"format"`
	Filters map[string]any `json:"filters"`
	Status  ExportStatus   `json:"status"`
	Rows    *int64         `json:"rows"`
	// The plan's export row limit cut the file short.
	Truncated  bool       `json:"truncated"`
	Bytes      *int64     `json:"bytes"`
	Error      *string    `json:"error"`
	CreatedAt  time.Time  `json:"createdAt"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
	// When the file is deleted (7 days after it was written).
	ExpiresAt time.Time `json:"expiresAt"`
	// Signed, credential-free URL of the file, valid 15 minutes from this response (set when succeeded).
	DownloadURL *string `json:"downloadUrl"`
}

// ExportPage is a schema of the public API.
type ExportPage struct {
	Items      []Export `json:"items"`
	NextCursor *string  `json:"nextCursor"`
}
