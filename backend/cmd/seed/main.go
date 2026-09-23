// Command seed populates a database with test users and a made-up chemical inventory across sample
// lab locations. The inventory data is entirely fabricated (see generateChemicalRows) so
// no real inventory data ever needs to be committed to this repo. It is meant for local/dev
// databases only — running it against a database that already has this data will create
// duplicate items (properties, item types and users are reused if they already exist).
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

const (
	// The manufacturer is deliberately absent: the derived name is also what stock is counted by, and
	// two bottles of the same chemical at the same purity are interchangeable however they were
	// bought. Naming the supplier here would split one stock line into one per supplier, so a shelf
	// holding plenty would read as several groups each nearly empty.
	chemicalDerivedNameFormat = "{Naziv hemikalije} · {Čistoća}"
	testPassword              = "Test1234"
	// chemicalExpiringSoonDays is how many days ahead of its date an item counts as expiring soon.
	// Set on the seeded item type, and reused to work out which items an expiry notification may
	// truthfully point at.
	chemicalExpiringSoonDays = 14
	// chemicalLowStockCount is how few available items of one derived name leave the type low on
	// stock. Set on the seeded item type so the seeded shelf trips the warning.
	chemicalLowStockCount = 2
)

// measure is a package size as it appears on the bottle: an amount in a named unit. It maps onto
// the mass/volume property types, which store the amount scaled by dto.MeasureMultiplier.
type measure struct {
	Amount float64
	Unit   string
}

type chemicalRow struct {
	Name string
	// CASNumber is the reagent's CAS registry number, the identifier a chemist actually searches
	// an inventory by. It identifies the substance, not this inventory, so it is real data.
	CASNumber    string
	Manufacturer string
	Purity       string
	// A package is stocked either by mass (solids) or by volume (liquids), never both, so exactly
	// one of these is set.
	PackageMass   *measure
	PackageVolume *measure
	// PricePerPackage is what one package cost, in RSD. Unlike mass/volume it is set on every row, so
	// the seed has spend data to sum regardless of which measure a chemical uses.
	PricePerPackage float64
	// PackageCount is how many identical packages of this reagent sit on the shelf. Rows with a
	// count above one are seeded through the bulk-add path in one call, the same way the UI stocks
	// several identical bottles at once, so every copy is its own item with its own request state.
	PackageCount int
	ExpiryDate   string
	// ExpiryOffsetDays is ExpiryDate expressed as days from today, kept alongside it so the expiry
	// state of an item can be told without parsing the date back out again.
	ExpiryOffsetDays int
}

// generateChemicalRows fabricates a chemical inventory shaped like a real one. Chemical names are
// real (common lab reagents), since a chemical's name isn't sensitive on its own; everything tied
// to this specific inventory (manufacturer, dates, box) is still made up so no real data is
// committed. Each row becomes PackageCount items. Selection is index-based (no randomness) so the
// output is stable across runs.
func generateChemicalRows() []chemicalRow {
	// solid marks reagents stocked as a powder or crystal, which a lab weighs out; everything else
	// is a liquid stocked by volume. The split is what gives the seed both mass and volume data.
	// Salts are listed under their anhydrous CAS number; the seed does not model hydrates, which
	// are separate registry entries.
	chemicals := []struct {
		name  string
		cas   string
		solid bool
	}{
		{"Acetone", "67-64-1", false},
		{"Methanol", "67-56-1", false},
		{"Ethanol", "64-17-5", false},
		{"Isopropanol", "67-63-0", false},
		{"Hexane", "110-54-3", false},
		{"Heptane", "142-82-5", false},
		{"Pentane", "109-66-0", false},
		{"Cyclohexane", "110-82-7", false},
		{"Toluene", "108-88-3", false},
		{"Xylene", "1330-20-7", false},
		{"Benzene", "71-43-2", false},
		{"Dichloromethane", "75-09-2", false},
		{"Chloroform", "67-66-3", false},
		{"Diethyl ether", "60-29-7", false},
		{"Tetrahydrofuran", "109-99-9", false},
		{"Dimethyl sulfoxide", "67-68-5", false},
		{"Acetonitrile", "75-05-8", false},
		{"Ethyl acetate", "141-78-6", false},
		{"Formaldehyde", "50-00-0", false},
		{"Glycerol", "56-81-5", false},
		{"Sulfuric acid", "7664-93-9", false},
		{"Hydrochloric acid", "7647-01-0", false},
		{"Nitric acid", "7697-37-2", false},
		{"Phosphoric acid", "7664-38-2", false},
		{"Acetic acid", "64-19-7", false},
		{"Oxalic acid", "144-62-7", true},
		{"Citric acid", "77-92-9", true},
		{"Tartaric acid", "87-69-4", true},
		{"Boric acid", "10043-35-3", true},
		{"Sodium hydroxide", "1310-73-2", true},
		{"Potassium hydroxide", "1310-58-3", true},
		{"Ammonium hydroxide", "1336-21-6", false},
		{"Sodium chloride", "7647-14-5", true},
		{"Potassium chloride", "7447-40-7", true},
		{"Sodium carbonate", "497-19-8", true},
		{"Sodium bicarbonate", "144-55-8", true},
		{"Sodium sulfate", "7757-82-6", true},
		{"Sodium acetate", "127-09-3", true},
		{"Sodium thiosulfate", "7772-98-7", true},
		{"Sodium hypochlorite", "7681-52-9", false},
		{"Ammonium chloride", "12125-02-9", true},
		{"Ammonium nitrate", "6484-52-2", true},
		{"Ammonium acetate", "631-61-8", true},
		{"Calcium chloride", "10043-52-4", true},
		{"Magnesium sulfate", "7487-88-9", true},
		{"Barium chloride", "10361-37-2", true},
		{"Zinc sulfate", "7733-02-0", true},
		{"Copper sulfate", "7758-98-7", true},
		{"Iron(III) chloride", "7705-08-0", true},
		{"Silver nitrate", "7761-88-8", true},
		{"Potassium iodide", "7681-11-0", true},
		{"Potassium permanganate", "7722-64-7", true},
		{"Hydrogen peroxide", "7722-84-1", false},
		{"Iodine", "7553-56-2", true},
		{"Bromine", "7726-95-6", false},
		{"EDTA", "60-00-4", true},
		{"Phenolphthalein", "77-09-8", true},
		{"Methyl orange", "547-58-0", true},
		{"Pyridine", "110-86-1", false},
		{"Triethylamine", "121-44-8", false},
	}
	manufacturers := []string{"NovaChem", "Solvex Labs", "Ferronova", "BluePeak Reagents", "Cryotech Supply", "Meridian Chemicals", "Vertex Labs", "Arcadia Chemical", "Lumen Scientific", "Pinegrove Labs"}
	purities := []string{"PA", "HPLC", "GC", "ultrapure", "technical", "0.99", "0.995", "0.997", "ACS"}
	massPackages := []measure{{1.0, "kg"}, {2.5, "kg"}, {500, "g"}, {100, "g"}, {25, "g"}}
	volumePackages := []measure{{1.0, "L"}, {2.5, "L"}, {5.0, "L"}, {500, "mL"}, {250, "mL"}}
	// Amounts of items to bulk add
	// 1 and 2 are there to trigger low stock threshold alerts
	packageCounts := []int{4, 3, 6, 3, 5, 4, 3, 1, 5, 3, 6, 2}
	// Days from today to each row's expiry date. Negatives are already expired and the small
	// positives fall inside the type's chemicalExpiringSoonDays window, so the shelf carries both
	// warning states rather than being uniformly fresh - without that, an expiry notification has
	// nothing truthful to point at. Cycled by index, so reruns produce the same inventory.
	expiryOffsetDays := []int{-45, 16, 120, 32, 40, 72, 200, 43, 70, 11, 21, 25, 9, 320, -7, 55}
	// Per-package price in RSD, cycled the same way as the other made-up fields: cheap solvents
	// through specialty reagents, in no particular order tied to the real chemical.
	packagePrices := []float64{950, 1450, 2200, 3100, 4200, 5800, 7900, 11000, 15500, 21000, 27500}

	rows := make([]chemicalRow, 0, len(chemicals))
	for i, chemical := range chemicals {
		expiryOffset := expiryOffsetDays[i%len(expiryOffsetDays)]
		expiryDate := time.Now().AddDate(0, 0, expiryOffset).Format(time.DateOnly)
		row := chemicalRow{
			Name:             chemical.name,
			CASNumber:        chemical.cas,
			Manufacturer:     manufacturers[i%len(manufacturers)],
			Purity:           purities[i%len(purities)],
			PackageCount:     packageCounts[i%len(packageCounts)],
			ExpiryDate:       expiryDate,
			ExpiryOffsetDays: expiryOffset,
			PricePerPackage:  packagePrices[i%len(packagePrices)],
		}
		if chemical.solid {
			packageMass := massPackages[i%len(massPackages)]
			row.PackageMass = &packageMass
		} else {
			packageVolume := volumePackages[i%len(volumePackages)]
			row.PackageVolume = &packageVolume
		}

		rows = append(rows, row)
	}

	return rows
}

type testUser struct {
	Email    string
	FullName string
	Role     string
	Status   string
}

var testUsers = []testUser{
	{Email: "admin@popisomator.test", FullName: "Admin Adminović", Role: "admin", Status: "active"},
	{Email: "manager@popisomator.test", FullName: "Manja Menadžer", Role: "manager", Status: "active"},
	{Email: "user1@popisomator.test", FullName: "Pera Perić", Role: "user", Status: "active"},
	{Email: "user2@popisomator.test", FullName: "Mika Mikić", Role: "user", Status: "active"},
	{Email: "jovana@popisomator.test", FullName: "Jovana Jovanović", Role: "user", Status: "requested"},
	{Email: "nikola@popisomator.test", FullName: "Nikola Nikolić", Role: "user", Status: "requested"},
}

// seededItem is a created item together with how far its expiry date sits from today. Carrying the
// offset out of seedItems means an expiry notification can be aimed at an item that genuinely is
// expired or near expiry, instead of at a hard-coded index that says nothing about its date.
type seededItem struct {
	Item             dto.Item
	ExpiryOffsetDays int
}

func (item seededItem) expiryState() (repository.NotifdescExpiryType, bool) {
	switch {
	case item.ExpiryOffsetDays < 0:
		return repository.NotifdescExpiryTypeExpired, true
	case item.ExpiryOffsetDays <= chemicalExpiringSoonDays:
		return repository.NotifdescExpiryTypeExpiringSoon, true
	default:
		return "", false
	}
}

type itemRequestSeed struct {
	Email string
	// ItemIndex points into the flat list of created items, not into generateChemicalRows: a
	// bulk-added reagent occupies as many consecutive slots as it has packages, so neighbouring
	// indexes can be two copies of the same chemical. That is intentional here - it leaves one
	// copy requested and the rest free, which is the case worth showing.
	ItemIndex int
	Reason    string
	Approved  bool
}

var itemRequestSeeds = []itemRequestSeed{
	{Email: "user1@popisomator.test", ItemIndex: 0, Reason: "Potreban mi je za pripremu rastvora."},
	{Email: "user2@popisomator.test", ItemIndex: 1},
	{Email: "user1@popisomator.test", ItemIndex: 2, Reason: "Za planirani eksperiment.", Approved: true},
	{Email: "user2@popisomator.test", ItemIndex: 3, Reason: "Dopuna zaliha za vežbe."},
	// Approved, and here so each test user has a decided request to be notified about; one approved
	// request between them would leave the other's notification list empty.
	{Email: "user1@popisomator.test", ItemIndex: 8, Reason: "Za analizu uzoraka.", Approved: true},
	{Email: "user2@popisomator.test", ItemIndex: 12, Reason: "Za redovne laboratorijske vežbe.", Approved: true},
}

// notificationSeed is one row in a recipient's notification list. Kind picks which descriptor the
// notification carries, and so which of the two index fields below is read.
type notificationSeed struct {
	// Email is the recipient - the person whose bell this lands on, not the person who caused it.
	Email string
	Kind  repository.NotificationKind
	// ItemRequestIndex points into itemRequestSeeds, and is only read for item_request. The
	// descriptor's foreign key is onto item_requests, so the request has to already exist.
	//
	// One kind covers two opposite messages, told apart only by whether the recipient is the person
	// who made the request. Someone else is being asked to decide it, so that request must still be
	// pending; the requester is being told it was decided, so theirs must be approved. Getting the
	// pairing backwards seeds a list that reads as nonsense - a manager chased about a request that
	// is already settled, or a user informed they requested their own item - so seedNotifications
	// rejects both rather than leaving it to be caught by eye.
	ItemRequestIndex int
	// ExpiryType says whether the item is merely approaching its date or already past it, and is
	// only read for item_expiry. There is deliberately no item index to go with it: seedNotifications
	// picks from the items whose dates actually put them in this state, so the item can never
	// contradict what the notification says about it.
	ExpiryType repository.NotifdescExpiryType
	// AgeHours backdates created_at. Real notifications trickle in over time and the list is
	// ordered newest first, so seeding every row at now() would leave nothing to sort by. Two seeds
	// sharing an age is deliberate: that is the identical-timestamp case the id tiebreaker in
	// ListNotifications exists to settle.
	AgeHours int
	// Read rows are the ones a recipient has already seen: no highlight, and sorted below every
	// unread row no matter how much older the unread one is.
	Read bool
}

// consumptionSeed is one item taken off the shelf.
type consumptionSeed struct {
	// ItemIndex points into the flat list of created items, not into the reagent rows.
	ItemIndex int
	State     string
	// ConsumedAt is set on the entry afterwards, since a real one is always stamped with the time it
	// was made.
	ConsumedAt time.Time
	// Only managers and admins, who may consume anything. A plain user is limited to items they hold
	// an approved request for.
	Email string
}

// consumptionsPerMonth is how many items were used up in each of the last 12 months, most recent
// first, weighted so the shortest range lands on the busy end of the year.
var consumptionsPerMonth = []int{6, 5, 5, 4, 4, 3, 3, 3, 2, 2, 2, 1}

// generateConsumptionSeeds decides which items have already been used, and when. It works from the
// end of the list so the items already spoken for by a seeded request keep the state they were
// given, and picks by index rather than at random so runs stay identical.
func generateConsumptionSeeds(itemCount int) []consumptionSeed {
	claimed := make(map[int]struct{}, len(itemRequestSeeds))
	for _, seed := range itemRequestSeeds {
		claimed[seed.ItemIndex] = struct{}{}
	}

	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 12, 0, 0, 0, time.UTC)

	seeds := make([]consumptionSeed, 0, itemCount)
	itemIndex := itemCount - 1
	for monthsAgo, count := range consumptionsPerMonth {
		monthStart := currentMonthStart.AddDate(0, -monthsAgo, 0)

		for slot := range count {
			for itemIndex >= 0 {
				if _, isClaimed := claimed[itemIndex]; !isClaimed {
					break
				}
				itemIndex--
			}
			if itemIndex < 0 {
				return seeds
			}

			// The current month is only as long as it has got so far, or the entry lands in the future.
			dayOfMonth := slot * 5
			if monthsAgo == 0 {
				dayOfMonth %= now.Day()
			} else {
				dayOfMonth %= 28
			}

			state := string(repository.ConsumptionStatusFullyConsumed)
			switch len(seeds) % 7 {
			case 5:
				state = string(repository.ConsumptionStatusPartiallyConsumed)
			case 6:
				state = string(repository.ConsumptionStatusDamaged)
			}

			email := "manager@popisomator.test"
			if len(seeds)%3 == 0 {
				email = "admin@popisomator.test"
			}

			seeds = append(seeds, consumptionSeed{
				ItemIndex:  itemIndex,
				State:      state,
				ConsumedAt: monthStart.AddDate(0, 0, dayOfMonth),
				Email:      email,
			})
			itemIndex--
		}
	}

	return seeds
}

// propertyDef describes one of the "Hemikalija" item type's properties.
type propertyDef struct {
	key        string // matches a key in propertyValues() below
	name       string
	valueType  string
	visibility repository.PropertyVisibility
}

var propertyDefs = []propertyDef{
	{"name", "Naziv hemikalije", "string", repository.PropertyVisibilityDetails},
	{"manufacturer", "Proizvođač", "string", repository.PropertyVisibilityDetails},
	{"cas_number", "CAS broj", "string", repository.PropertyVisibilityOverview},
	{"purity", "Čistoća", "string", repository.PropertyVisibilityOverview},
	{"mass", "Masa", "mass", repository.PropertyVisibilityOverview},
	{"volume", "Zapremina", "volume", repository.PropertyVisibilityOverview},
	{"price", "Cena", "price", repository.PropertyVisibilityOverview},
	{"expiry_date", "Istek roka", "expiry", repository.PropertyVisibilityOverview},
}

type locationSeed struct {
	name        string
	description string
	children    []locationSeed
}

var locationSeeds = []locationSeed{
	{
		name:        "Laboratorija 1",
		description: "Glavni radni prostor sa pristupom ventilaciji i centralnim radnim pultovima.",
		children: []locationSeed{
			{
				name:        "Sigurnosni ormar A",
				description: "Vatrostalni ormar sa nezavisnim sistemom za ventilaciju.",
				children: []locationSeed{
					{
						name:        "Polica A-1",
						description: "Gornja polica, maksimalna nosivost 20kg.",
					},
					{
						name:        "Polica A-2",
						description: "Donja polica, ojačana, opremljena sabirnom kadicom za slučaj izlivanja.",
						children: []locationSeed{
							{
								name:        "Zaštitna kutija A-2-1",
								description: "Prenosiva polipropilenska kutija, koristi se kao sekundarni kontejner.",
							},
							{
								name:        "Zaštitna kutija A-2-2",
								description: "Prenosiva polipropilenska kutija, koristi se kao sekundarni kontejner.",
							},
						},
					},
				},
			},
			{
				name:        "Zidni ormar B",
				description: "Metalni ormar sa staklenim vratima, koristi se za skladištenje na sobnoj temperaturi.",
				children: []locationSeed{
					{
						name:        "Odeljak B-Levi",
						description: "Levi deo ormara sa tri podesive police.",
					},
					{
						name:        "Odeljak B-Desni",
						description: "Desni deo ormara sa fiksnom pregradom.",
					},
				},
			},
		},
	},
	{
		name:        "Centralni magacin",
		description: "Prostorija za dugoročno skladištenje sa kontrolisanim pristupom i temperaturom.",
		children: []locationSeed{
			{
				name:        "Frižider F-1",
				description: "Laboratorijski frižider sa konstantnim temperaturnim režimom (2-8°C).",
				children: []locationSeed{
					{
						name:        "Fioka 1",
						description: "Gornja plitka fioka od nerđajućeg čelika.",
					},
					{
						name:        "Fioka 2",
						description: "Donja duboka fioka od nerđajućeg čelika.",
					},
				},
			},
			{
				name:        "Stalaža S-1",
				description: "Industrijska metalna stalaža otvorenog tipa.",
				children: []locationSeed{
					{
						name:        "Nivo 1 (Podni)",
						description: "Prostor namenjen za tešku ambalažu, direktno iznad poda.",
					},
					{
						name:        "Nivo 2 (Srednji)",
						description: "Središnja polica, standardna visina za lako pristupanje.",
						children: []locationSeed{
							{
								name:        "Plastični kontejner S-1-A",
								description: "Otvoreni plastični kontejner za organizaciju manjih boca.",
							},
						},
					},
				},
			},
		},
	},
}

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("unable to initialise config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer pool.Close()

	users, err := seedUsers(ctx)
	if err != nil {
		log.Fatalf("unable to seed users: %v", err)
	}

	// Everything below runs through the same services the API does, so each call records its own
	// audit entry. Naming an actor is what keeps those entries from all reading "Sistem": the
	// catalog and the stock are the admin's doing, the same way they would be in the running app.
	admin, ok := users["admin@popisomator.test"]
	if !ok {
		log.Fatal("seed admin was not created")
	}
	adminCtx := actorContext(ctx, admin.ID)

	propIDs, err := seedProperties(adminCtx)
	if err != nil {
		log.Fatalf("unable to seed properties: %v", err)
	}

	typeID, err := seedItemType(adminCtx, "Hemikalija", propIDs)
	if err != nil {
		log.Fatalf("unable to seed item type: %v", err)
	}

	locationIDs, err := seedLocations(ctx)
	if err != nil {
		log.Fatalf("unable to seed locations: %v", err)
	}

	items, err := seedItems(adminCtx, typeID, propIDs, locationIDs)
	if err != nil {
		log.Fatalf("unable to seed items: %v", err)
	}

	if err := seedItemRequests(ctx, adminCtx, users, items); err != nil {
		log.Fatalf("unable to seed item requests: %v", err)
	}

	if err := seedConsumption(ctx, users, items); err != nil {
		log.Fatalf("unable to seed consumption: %v", err)
	}

	if err := trimLowStockNotifications(ctx); err != nil {
		log.Fatalf("unable to trim low stock notifications: %v", err)
	}

	fmt.Println("seeding complete")
}

// actorContext attributes whatever the services do with it to userID, under the same context key
// middleware.RequireAuth uses on a real request.
func actorContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func seedUsers(ctx context.Context) (map[string]dto.User, error) {
	users := make(map[string]dto.User, len(testUsers))

	for _, testUser := range testUsers {
		existing, err := service.GetUserByEmail(ctx, testUser.Email)
		if err == nil {
			fmt.Printf("user %s already exists, skipping\n", testUser.Email)
			users[testUser.Email] = existing
			continue
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		created, err := service.CreateUser(ctx, dto.CreateUserRequest{
			Email:    testUser.Email,
			Password: testPassword,
			FullName: testUser.FullName,
			Role:     testUser.Role,
			Status:   testUser.Status,
		})
		if err != nil {
			return nil, fmt.Errorf("creating user %s: %w", testUser.Email, err)
		}
		fmt.Printf("created user %d (%s / %s)\n", created.ID, created.Email, testPassword)
		users[testUser.Email] = created
	}

	return users, nil
}

func seedProperties(ctx context.Context) (map[string]int64, error) {
	existing, err := service.GetAllProperties(ctx)
	if err != nil {
		return nil, err
	}
	byName := make(map[string]int64, len(existing))
	for _, property := range existing {
		byName[property.Name] = property.ID
	}

	propIDs := make(map[string]int64, len(propertyDefs))
	for _, def := range propertyDefs {
		if propertyID, ok := byName[def.name]; ok {
			propIDs[def.key] = propertyID
			continue
		}

		created, err := service.CreateProperty(ctx, dto.CreatePropertyRequest{
			Name:      def.name,
			ValueType: def.valueType,
		})
		if err != nil {
			return nil, fmt.Errorf("creating property %s: %w", def.name, err)
		}
		fmt.Printf("created property %d (%s)\n", created.ID, created.Name)
		propIDs[def.key] = created.ID
	}

	return propIDs, nil
}

func seedItemType(ctx context.Context, name string, propIDs map[string]int64) (int64, error) {
	existing, err := service.GetAllItemTypes(ctx)
	if err != nil {
		return 0, err
	}
	for _, existingType := range existing {
		if existingType.Name == name {
			fmt.Printf("item type %s already exists (%d), ensuring properties are attached\n", name, existingType.ID)
			existingType, err := service.GetItemType(ctx, existingType.ID)
			if err != nil {
				return 0, fmt.Errorf("loading item type properties: %w", err)
			}
			if err := ensureItemTypeProperties(ctx, existingType, propIDs); err != nil {
				return 0, err
			}
			lowStockCount := int32(chemicalLowStockCount)
			if existingType.DerivedNameFormat != chemicalDerivedNameFormat ||
				existingType.LowStockCount == nil || *existingType.LowStockCount != lowStockCount {
				format := chemicalDerivedNameFormat
				if _, err := service.UpdateItemType(ctx, dto.UpdateItemTypeRequest{
					ID:                existingType.ID,
					DerivedNameFormat: &format,
					LowStockCount:     &lowStockCount,
				}); err != nil {
					return 0, fmt.Errorf("setting derived name format: %w", err)
				}
			}
			return existingType.ID, nil
		}
	}

	properties := make([]dto.ItemTypeProperty, 0, len(propertyDefs))
	for _, def := range propertyDefs {
		properties = append(properties, dto.ItemTypeProperty{
			ID:         propIDs[def.key],
			Visibility: string(def.visibility),
		})
	}

	expiringSoonDays := int16(chemicalExpiringSoonDays)
	lowStockCount := int32(chemicalLowStockCount)

	created, err := service.CreateItemType(ctx, dto.CreateItemTypeRequest{
		Name:              name,
		DerivedNameFormat: chemicalDerivedNameFormat,
		Properties:        properties,
		ExpiringSoonDays:  &expiringSoonDays,
		LowStockCount:     &lowStockCount,
	})
	if err != nil {
		return 0, err
	}
	fmt.Printf("created item type %d (%s)\n", created.ID, created.Name)

	return created.ID, nil
}

func ensureItemTypeProperties(ctx context.Context, itemType dto.ItemType, propIDs map[string]int64) error {
	attached := make(map[int64]dto.ItemTypeProperty, len(itemType.Properties))
	for _, property := range itemType.Properties {
		attached[property.ID] = property
	}

	for _, def := range propertyDefs {
		propertyID := propIDs[def.key]
		visibility := def.visibility
		property, exists := attached[propertyID]
		if !exists {
			if _, err := service.AddItemTypeProperty(ctx, dto.AddUpdateItemTypePropertyRequest{
				TypeID:     itemType.ID,
				PropertyID: propertyID,
				Visibility: &visibility,
			}); err != nil {
				return fmt.Errorf("attaching property %s to item type: %w", def.name, err)
			}
			continue
		}

		if property.Visibility != string(visibility) {
			if _, err := service.UpdateItemTypeProperty(ctx, dto.AddUpdateItemTypePropertyRequest{
				TypeID:     itemType.ID,
				PropertyID: propertyID,
				Visibility: &visibility,
			}); err != nil {
				return fmt.Errorf("setting visibility for property %s: %w", def.name, err)
			}
		}
	}

	return nil
}

func createLocation(ctx context.Context, entity locationSeed, parentID *int64) ([]int64, error) {
	location, err := service.CreateLocation(ctx, dto.CreateLocationRequest{
		Name:        entity.name,
		Description: entity.description,
		ParentID:    parentID,
	})
	if err != nil {
		return nil, err
	}
	IDs := make([]int64, 0, 1)
	IDs = append(IDs, location.ID)

	if entity.children != nil {
		for _, child := range entity.children {
			childIDs, err := createLocation(ctx, child, &location.ID)
			if err != nil {
				return nil, err
			}
			IDs = append(IDs, childIDs...)
		}
	}

	return IDs, nil
}

func seedLocations(ctx context.Context) ([]int64, error) {
	locationIDs := make([]int64, 0)
	for _, entity := range locationSeeds {
		entityIDs, err := createLocation(ctx, entity, nil)
		if err != nil {
			return nil, err
		}
		locationIDs = append(locationIDs, entityIDs...)
	}

	return locationIDs, nil
}

// seedItems creates one item per physical package rather than one item per spreadsheet row, so a
// row stocked in several identical packages becomes that many items. Those rows go through the
// same bulk-add path the UI uses (one CreateItem call with Amount > 1), which is what gives the
// seeded inventory duplicate items to exercise filtering and the property totals against.
func seedItems(ctx context.Context, typeID int64, propIDs map[string]int64, locationIDs []int64) ([]seededItem, error) {
	rows := generateChemicalRows()
	items := make([]seededItem, 0, len(rows))

	bulkRows := 0
	for i, row := range rows {
		properties := propertyValues(row, propIDs)

		createdItems, err := service.CreateItem(ctx, dto.CreateItemRequest{
			TypeID:     typeID,
			Properties: properties,
			Amount:     int32(row.PackageCount),
			LocationID: &locationIDs[i%len(locationIDs)],
		})
		if err != nil {
			return nil, fmt.Errorf("creating item %q: %w", row.Name, err)
		}

		itemIDs := make([]string, len(createdItems))
		for index, createdItem := range createdItems {
			itemIDs[index] = strconv.FormatInt(createdItem.ID, 10)
		}
		fmt.Printf("created %d item(s) %s (%s)\n", len(createdItems), strings.Join(itemIDs, ", "), row.Name)

		for _, createdItem := range createdItems {
			items = append(items, seededItem{Item: createdItem, ExpiryOffsetDays: row.ExpiryOffsetDays})
		}
		if len(createdItems) > 1 {
			bulkRows++
		}
	}

	fmt.Printf("seeded %d items across sample locations (%d of %d reagents bulk-added)\n",
		len(items), bulkRows, len(rows))
	return items, nil
}

// seedItemRequests takes two contexts because a request and its approval are two different people's
// doing: the user asks for the item, the admin grants it. Recording them under one actor would make
// the resulting audit log a poor sample of the real thing.
func seedItemRequests(ctx, adminCtx context.Context, users map[string]dto.User, items []seededItem) error {
	for _, seed := range itemRequestSeeds {
		user, ok := users[seed.Email]
		if !ok {
			return fmt.Errorf("seed user %s was not created", seed.Email)
		}
		if seed.ItemIndex >= len(items) {
			return fmt.Errorf("item request references missing item at index %d", seed.ItemIndex)
		}

		item := items[seed.ItemIndex].Item
		itemRequest, err := service.CreateItemRequest(actorContext(ctx, user.ID), dto.ItemRequestCreateRequest{
			UserID: user.ID,
			ItemID: item.ID,
			Reason: seed.Reason,
		})
		if err != nil {
			return fmt.Errorf("creating request for item %d: %w", item.ID, err)
		}

		requestStatus := "pending"
		if seed.Approved {
			if _, err := service.ApproveItemRequest(adminCtx, dto.ItemRequestIdentifierRequest{
				UserID: itemRequest.UserID,
				ItemID: itemRequest.ItemID,
			}); err != nil {
				return fmt.Errorf("approving request for item %d: %w", item.ID, err)
			}
			requestStatus = "approved"
		}

		fmt.Printf("created %s request for item %d (%s)\n", requestStatus, item.ID, user.FullName)
	}

	return nil
}

// seedConsumption uses up part of the seeded shelf. It goes through the service rather than writing
// the status directly, so each one leaves a real entry behind it - which is where every usage figure
// is read from, and nothing else in the seed produces any.
func seedConsumption(ctx context.Context, users map[string]dto.User, items []seededItem) error {
	seeds := generateConsumptionSeeds(len(items))

	itemIDs := make([]int64, 0, len(seeds))
	consumedAts := make([]time.Time, 0, len(seeds))

	for _, seed := range seeds {
		actor, ok := users[seed.Email]
		if !ok {
			return fmt.Errorf("seed user %s was not created", seed.Email)
		}
		if seed.ItemIndex >= len(items) {
			return fmt.Errorf("consumption seed points at item %d of %d", seed.ItemIndex, len(items))
		}
		item := items[seed.ItemIndex].Item

		state := seed.State
		if _, err := service.UpdateItem(actorContext(ctx, actor.ID), dto.UpdateItemRequest{
			ID:          item.ID,
			Consumption: &state,
			ViewerID:    actor.ID,
		}); err != nil {
			return fmt.Errorf("consuming item %d: %w", item.ID, err)
		}

		itemIDs = append(itemIDs, item.ID)
		consumedAts = append(consumedAts, seed.ConsumedAt)
	}

	fmt.Printf("consumed %d of %d items across the last %d months\n",
		len(seeds), len(items), len(consumptionsPerMonth))

	return backdateConsumption(ctx, itemIDs, consumedAts)
}

// backdateConsumption moves each seeded consumption, and the item it happened to, back to the date
// it should read as. Both are stamped with the current time on insert, which is right everywhere
// except in a seed that needs a year of history behind it - so the statements stay here rather than
// going into the shared query set.
//
// The item's own creation moves with it: nothing can be used before it was stocked, and a timeline
// that opens with the item being consumed looks like a bug rather than like old data. Items nothing
// was seeded against are left alone.
func backdateConsumption(ctx context.Context, itemIDs []int64, consumedAts []time.Time) error {
	if len(itemIDs) == 0 {
		return nil
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// A stocking date of one to three months before the item was used, varied by id so the items do
	// not all arrive together.
	if _, err := tx.Exec(ctx, `
		UPDATE items AS item
		SET created_at = seed.consumed_at - make_interval(days => 30 + (item.id % 60)::int)
		FROM (
			SELECT
				unnest($1::bigint[]) AS id,
				unnest($2::timestamptz[]) AS consumed_at
		) AS seed
		WHERE item.id = seed.id
	`, itemIDs, consumedAts); err != nil {
		return fmt.Errorf("backdating consumed items: %w", err)
	}

	// Matched on the item rather than on the entry, whose id a bulk add never hands back.
	if _, err := tx.Exec(ctx, `
		UPDATE audit_log AS entry
		SET created_at = CASE entry.action
			WHEN 'item_consume' THEN seed.consumed_at
			ELSE (SELECT items.created_at FROM items WHERE items.id = seed.id)
		END
		FROM (
			SELECT
				unnest($1::bigint[]) AS id,
				unnest($2::timestamptz[]) AS consumed_at
		) AS seed
		WHERE entry.target_type = 'item'
		  AND entry.target_id = seed.id
		  AND entry.action IN ('item_create', 'item_consume')
	`, itemIDs, consumedAts); err != nil {
		return fmt.Errorf("backdating consumption audit entries: %w", err)
	}

	return tx.Commit(ctx)
}

// keptLowStockNotifications is how many of the warnings below survive per recipient - enough that
// the notification list has the kind in it, few enough that it does not drown out the rest.
const keptLowStockNotifications = 3

// trimLowStockNotifications throws away most of the low stock warnings the seed sets off.
//
// Stocking the shelf and using part of it both run through the services, which tell every manager
// and admin each time a group falls to its threshold. Most reagents here come in ones and twos
// against a threshold of two, so that is most of the catalogue warned about at once, and a hundred
// near-identical rows bury the handful of notifications this seed actually composed.
//
// The alert rows are deliberately left alone. They are what records a group as already warned
// about, so the shortages stay known and nothing fires again the moment the app is touched - this
// only clears the inbox, exactly as it would look if someone had read and dismissed them.
func trimLowStockNotifications(ctx context.Context) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		DELETE FROM notifications
		WHERE kind = 'item_low_stock'
		  AND id NOT IN (
			SELECT id FROM (
				SELECT id, row_number() OVER (PARTITION BY recipient_id ORDER BY id DESC) AS position
				FROM notifications
				WHERE kind = 'item_low_stock'
			) AS ranked
			WHERE position <= $1
		)
	`, keptLowStockNotifications)
	if err != nil {
		return fmt.Errorf("trimming low stock notifications: %w", err)
	}

	fmt.Printf("dropped %d low stock notifications raised while seeding, kept %d per recipient\n",
		tag.RowsAffected(), keptLowStockNotifications)

	return tx.Commit(ctx)
}

// propertyValues builds the property list for a row's item, skipping any field left blank
// rather than writing empty/zero values.
func propertyValues(row chemicalRow, propIDs map[string]int64) []dto.ItemProperty {
	properties := make([]dto.ItemProperty, 0, len(propertyDefs))

	addString := func(key, value string) {
		if value == "" {
			return
		}
		properties = append(properties, dto.ItemProperty{ID: propIDs[key], Value: json.RawMessage(fmt.Sprintf("\"%v\"", value))})
	}
	addMeasure := func(key string, value *measure, encode func(measure) any) {
		if value == nil {
			return
		}
		raw, err := json.Marshal(encode(*value))
		if err != nil {
			return
		}
		properties = append(properties, dto.ItemProperty{ID: propIDs[key], Value: raw})
	}

	addString("name", row.Name)
	addString("manufacturer", row.Manufacturer)
	addString("cas_number", row.CASNumber)
	addString("purity", formatPurity(row.Purity))
	addMeasure("mass", row.PackageMass, func(packageMass measure) any {
		return dto.PTMass{Amount: scaleMeasureAmount(packageMass.Amount), Unit: packageMass.Unit}
	})
	addMeasure("volume", row.PackageVolume, func(packageVolume measure) any {
		return dto.PTVolume{Amount: scaleMeasureAmount(packageVolume.Amount), Unit: packageVolume.Unit}
	})
	if row.PricePerPackage > 0 {
		raw, err := json.Marshal(dto.PTPrice{
			Amount:   int64(math.Round(row.PricePerPackage * dto.PriceMultiplier)),
			Currency: "RSD",
		})
		if err == nil {
			properties = append(properties, dto.ItemProperty{ID: propIDs["price"], Value: raw})
		}
	}
	addString("expiry_date", row.ExpiryDate)

	return properties
}

// scaleMeasureAmount converts a package size as written on the bottle into the integer amount the
// mass/volume property types store, i.e. scaled by dto.MeasureMultiplier.
func scaleMeasureAmount(amount float64) int64 {
	return int64(math.Round(amount * dto.MeasureMultiplier))
}

// formatPurity turns a bare decimal fraction (e.g. "0.995", as some rows record purity) into a
// percentage string (e.g. "99.5%"). Values that aren't a plain fraction (grades like "P.A." or
// "HPLC", or already-percented strings) pass through unchanged.
func formatPurity(purity string) string {
	fraction, err := strconv.ParseFloat(purity, 64)
	if err != nil {
		return purity
	}
	return strconv.FormatFloat(fraction*100, 'f', -1, 64) + "%"
}
