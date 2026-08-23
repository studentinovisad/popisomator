// Command seed populates a database with test users and a made-up chemical inventory for a
// single lab location. The inventory data is entirely fabricated (see generateChemicalRows) so
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
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

const (
	location                  = "Laboratorija A"
	chemicalDerivedNameFormat = "{Naziv hemikalije} · {Proizvođač}"
	testPassword              = "Test1234"
)

type chemicalRow struct {
	Name           string
	Manufacturer   string
	Purity         string
	PackageAmount  *float64
	PackageUnit    string
	PersonInCharge string
	NoteDate       string
	Box            string
}

// generateChemicalRows fabricates a chemical inventory shaped like a real one. Chemical names are
// real (common lab reagents), since a chemical's name isn't sensitive on its own; everything tied
// to this specific inventory (manufacturer, person in charge, dates, box) is still made up so no
// real data is committed. Each row becomes exactly one item. Selection is index-based (no
// randomness) so the output is stable across runs.
func generateChemicalRows() []chemicalRow {
	chemicalNames := []string{
		"Acetone", "Methanol", "Ethanol", "Isopropanol", "Hexane", "Heptane", "Pentane",
		"Cyclohexane", "Toluene", "Xylene", "Benzene", "Dichloromethane", "Chloroform",
		"Diethyl ether", "Tetrahydrofuran", "Dimethyl sulfoxide", "Acetonitrile", "Ethyl acetate",
		"Formaldehyde", "Glycerol", "Sulfuric acid", "Hydrochloric acid", "Nitric acid",
		"Phosphoric acid", "Acetic acid", "Oxalic acid", "Citric acid", "Tartaric acid",
		"Boric acid", "Sodium hydroxide", "Potassium hydroxide", "Ammonium hydroxide",
		"Sodium chloride", "Potassium chloride", "Sodium carbonate", "Sodium bicarbonate",
		"Sodium sulfate", "Sodium acetate", "Sodium thiosulfate", "Sodium hypochlorite",
		"Ammonium chloride", "Ammonium nitrate", "Ammonium acetate", "Calcium chloride",
		"Magnesium sulfate", "Barium chloride", "Zinc sulfate", "Copper sulfate",
		"Iron(III) chloride", "Silver nitrate", "Potassium iodide", "Potassium permanganate",
		"Hydrogen peroxide", "Iodine", "Bromine", "EDTA", "Phenolphthalein", "Methyl orange",
		"Pyridine", "Triethylamine",
	}
	manufacturers := []string{"NovaChem", "Solvex Labs", "Ferronova", "BluePeak Reagents", "Cryotech Supply", "Meridian Chemicals", "Vertex Labs", "Arcadia Chemical", "Lumen Scientific", "Pinegrove Labs"}
	purities := []string{"PA", "HPLC", "GC", "ultrapure", "technical", "0.99", "0.995", "0.997", "ACS"}
	packageAmounts := []float64{1.0, 2.5, 5.0}
	persons := []string{"A.B.", "C.D.", "M.N.", "J.K.", "T.R.", "S.P."}
	noteDates := []string{"", "", "12.03.2023", "05.07.2024", "21.11.2022", ""}
	boxes := []string{"", "", "Box 1", "Box 2", "Box 3", ""}

	rowCount := len(chemicalNames)
	rows := make([]chemicalRow, 0, rowCount)
	for i := 0; i < rowCount; i++ {
		packageAmount := packageAmounts[i%len(packageAmounts)]

		rows = append(rows, chemicalRow{
			Name:           chemicalNames[i],
			Manufacturer:   manufacturers[i%len(manufacturers)],
			Purity:         purities[i%len(purities)],
			PackageAmount:  &packageAmount,
			PackageUnit:    "L",
			PersonInCharge: persons[i%len(persons)],
			NoteDate:       noteDates[i%len(noteDates)],
			Box:            boxes[i%len(boxes)],
		})
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

type itemRequestSeed struct {
	Email     string
	ItemIndex int
	Reason    string
	Approved  bool
}

var itemRequestSeeds = []itemRequestSeed{
	{Email: "user1@popisomator.test", ItemIndex: 0, Reason: "Potreban mi je za pripremu rastvora."},
	{Email: "user2@popisomator.test", ItemIndex: 1},
	{Email: "user1@popisomator.test", ItemIndex: 2, Reason: "Za planirani eksperiment.", Approved: true},
	{Email: "user2@popisomator.test", ItemIndex: 3, Reason: "Dopuna zaliha za vežbe."},
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
	{"purity", "Čistoća", "string", repository.PropertyVisibilityOverview},
	{"total_amount", "Zapremina", "number", repository.PropertyVisibilityOverview},
	{"total_unit", "Jedinica mere", "string", repository.PropertyVisibilityOverview},
	{"person_in_charge", "Zadužena osoba", "string", repository.PropertyVisibilityDetails},
	{"note_date", "Napomena/datum", "string", repository.PropertyVisibilityDetails},
	{"box", "Mesto/kutija", "string", repository.PropertyVisibilityOverview},
	{"location", "Lokacija", "string", repository.PropertyVisibilityOverview},
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

	propIDs, err := seedProperties(ctx)
	if err != nil {
		log.Fatalf("unable to seed properties: %v", err)
	}

	typeID, err := seedItemType(ctx, "Hemikalija", propIDs)
	if err != nil {
		log.Fatalf("unable to seed item type: %v", err)
	}

	items, err := seedItems(ctx, typeID, propIDs)
	if err != nil {
		log.Fatalf("unable to seed items: %v", err)
	}

	if err := seedItemRequests(ctx, users, items); err != nil {
		log.Fatalf("unable to seed item requests: %v", err)
	}

	fmt.Println("seeding complete")
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
			if existingType.DerivedNameFormat != chemicalDerivedNameFormat {
				format := chemicalDerivedNameFormat
				if _, err := service.UpdateItemType(ctx, dto.UpdateItemTypeRequest{
					ID:                existingType.ID,
					DerivedNameFormat: &format,
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

	created, err := service.CreateItemType(ctx, dto.CreateItemTypeRequest{
		Name:              name,
		DerivedNameFormat: chemicalDerivedNameFormat,
		Properties:        properties,
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

// seedItems creates one item per physical package rather than one item per spreadsheet row: a
// each row becomes exactly one item.
func seedItems(ctx context.Context, typeID int64, propIDs map[string]int64) ([]dto.Item, error) {
	rows := generateChemicalRows()
	items := make([]dto.Item, 0, len(rows))

	created := 0
	for _, row := range rows {
		properties := propertyValues(row, propIDs)

		createdItems, err := service.CreateItem(ctx, dto.CreateItemRequest{
			TypeID:     typeID,
			Properties: properties,
			Amount:     1,
		})
		if err != nil {
			return nil, fmt.Errorf("creating item %q: %w", row.Name, err)
		}
		itemCreated := createdItems[0]
		fmt.Printf("created item %d (%s)\n", itemCreated.ID, row.Name)
		items = append(items, itemCreated)
		created++
	}

	fmt.Printf("seeded %d items for location %q\n", created, location)
	return items, nil
}

func seedItemRequests(ctx context.Context, users map[string]dto.User, items []dto.Item) error {
	for _, seed := range itemRequestSeeds {
		user, ok := users[seed.Email]
		if !ok {
			return fmt.Errorf("seed user %s was not created", seed.Email)
		}
		if seed.ItemIndex >= len(items) {
			return fmt.Errorf("item request references missing item at index %d", seed.ItemIndex)
		}

		item := items[seed.ItemIndex]
		itemRequest, err := service.CreateItemRequest(ctx, dto.ItemRequestCreateRequest{
			UserID: user.ID,
			ItemID: item.ID,
			Reason: seed.Reason,
		})
		if err != nil {
			return fmt.Errorf("creating request for item %d: %w", item.ID, err)
		}

		requestStatus := "pending"
		if seed.Approved {
			if _, err := service.ApproveItemRequest(ctx, dto.ItemRequestIdentifierRequest{
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

// propertyValues builds the property list for a row's item, skipping any field left blank
// rather than writing empty/zero values.
func propertyValues(row chemicalRow, propIDs map[string]int64) []dto.ItemProperty {
	properties := make([]dto.ItemProperty, 0, len(propertyDefs))

	addString := func(key, value string) {
		if value == "" {
			return
		}
		properties = append(properties, dto.ItemProperty{ID: propIDs[key], Value: jsonString(value)})
	}
	addNumber := func(key string, value *float64) {
		if value == nil {
			return
		}
		properties = append(properties, dto.ItemProperty{ID: propIDs[key], Value: jsonNumber(*value)})
	}

	addString("name", row.Name)
	addString("manufacturer", row.Manufacturer)
	addString("purity", formatPurity(row.Purity))
	addNumber("total_amount", row.PackageAmount)
	addString("total_unit", row.PackageUnit)
	addString("person_in_charge", row.PersonInCharge)
	addString("note_date", row.NoteDate)
	addString("box", row.Box)
	addString("location", location)

	return properties
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

func jsonString(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func jsonNumber(value float64) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}
