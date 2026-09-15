export type PageRequest = {
	limit?: number;
	offset?: number;
	search?: string;
};

export type UserRole = 'admin' | 'manager' | 'user';
export type UserStatus = 'requested' | 'active';

export type User = {
	id: number;
	email: string;
	full_name: string;
	role: UserRole;
	status: UserStatus;
};

export type UsersPage = {
	items: User[];
	limit: number;
	offset: number;
	total: number;
};

export type ListUsersParams = PageRequest & {
	search?: string;
	role?: UserRole;
	status?: UserStatus;
};

export type LoginRequest = {
	email: string;
	password: string;
};

export type RegistrationRequest = LoginRequest & {
	full_name: string;
};

export type CreateUserRequest = RegistrationRequest & {
	role: UserRole;
};

export type UpdateUserRequest = {
	email?: string;
	full_name?: string;
	role?: UserRole;
	status?: UserStatus;
};

export type ConsumptionStatus =
	'not_consumed' | 'partially_consumed' | 'fully_consumed' | 'damaged';

export type PropertyVisibility = 'overview' | 'details';

// Property values are primitives or JSON-like structured values. `object` is intentional here:
// individual value types are validated by the backend according to `PropertyValueType`.
export type PropertyValue = string | number | boolean | object;

export type ItemProperty = {
	id: number;
	value: PropertyValue;
	value_type?: PropertyValueType;
	visibility?: PropertyVisibility;
	smart_data?: string;
};

export type ItemRequestStatus = 'requested' | 'approved';

export type Item = {
	id: number;
	consumption: ConsumptionStatus;
	properties: ItemProperty[];
	type_id: number;
	derived_name?: string;
	request_status?: ItemRequestStatus;
	holder_name?: string;
	location_id?: number;
};

export type ItemsPage = {
	items: Item[];
	limit: number;
	offset: number;
	total: number;
	// Sums of the structured properties across every item matching the filters, not just this page.
	property_totals: ItemPropertyTotal[];
};

// Sum of one structured property (price, mass, volume) across every item matching a filter. `value`
// has the same shape as the property itself - a PTPrice, or a PTMeasure already converted to a
// readable unit - so displayJson renders it like any other value. Prices are reported per currency,
// so one property can appear more than once.
export type ItemPropertyTotal = {
	property_id: number;
	value_type: PropertyValueType;
	value: object;
	// How many of the matched items carried the property; compare with ItemsPage.total.
	value_count: number;
};

export type SortOrder = 'asc' | 'desc';
export type HeldBy = 'nobody' | 'me';

export type ListItemsParams = PageRequest & {
	typeID?: number;
	locationID?: number;
	propertyFilters?: Record<number, PropertyValue>;
	// Property to order the items by. Left out, they come back newest first.
	sortPropertyID?: number;
	order?: SortOrder;
	heldBy?: HeldBy;
};

export type UpdateItemRequest = {
	type_id?: number;
	location_id?: number | null;
	consumption?: ConsumptionStatus;
};

export type CreateItemRequest = {
	type_id: number;
	location_id?: number;
	properties: ItemProperty[];
	amount: number;
};

export type ItemTypeProperty = {
	id: number;
	default_value: PropertyValue | null;
	name?: string;
	visibility: PropertyVisibility;
};

export type ItemTypeFilterableProperty = {
	property_id: number;
	value_count: number;
};

export type ItemTypeOption = {
	id: number;
	name: string;
};

export type ItemType = ItemTypeOption & {
	description: string;
	properties: ItemTypeProperty[];
	derived_name_format: string;
	// How many days ahead of its date an expiry on this type counts as expiring soon. Null where no
	// window is set, which is how the backend records 'never mark anything as expiring soon'.
	expiring_soon_days: number | null;
	// How few available items sharing one derived name leave this type low on stock. Null where no
	// threshold is set, recorded the same way as expiring_soon_days.
	low_stock_count: number | null;
};

export type ItemTypesPage = {
	items: ItemType[];
	limit: number;
	offset: number;
	total: number;
};

export type CreateItemTypeRequest = {
	name: string;
	description: string;
	derived_name_format: string;
	properties: ItemTypeProperty[];
	expiring_soon_days?: number;
	low_stock_count?: number;
};

export type UpdateItemTypeRequest = Partial<
	Pick<
		ItemType,
		'name' | 'description' | 'derived_name_format' | 'expiring_soon_days' | 'low_stock_count'
	>
>;

export type AddUpdateItemTypePropertyRequest = {
	property_id?: number;
	default_value?: PropertyValue;
	visibility?: string;
};

export type PropertyValueType =
	'string' | 'number' | 'boolean' | 'price' | 'expiry' | 'mass' | 'volume';

export type PTPrice = {
	amount: number;
	currency: string;
};

// Shape shared by the measured property types (mass, volume): an amount in the unit the user
// picked, scaled by MeasureMultiplier.
export type PTMeasure = {
	amount: number;
	unit: string;
};

export type PropertyOption = {
	id: number;
	name: string;
	value_type: PropertyValueType;
	default_value: PropertyValue | null;
};

export type Property = PropertyOption & {
	description: string;
};

export type PropertiesPage = {
	items: Property[];
	limit: number;
	offset: number;
	total: number;
};

export type CreatePropertyRequest = {
	name: string;
	description: string;
	value_type: PropertyValueType;
	default_value: PropertyValue | null;
};

export type UpdatePropertyRequest = Partial<
	Pick<CreatePropertyRequest, 'name' | 'description' | 'default_value'>
>;

export type Location = {
	id: number;
	name: string;
	description?: string;
	parent_id?: number;
};

export type LocationOption = Pick<Location, 'id' | 'name'> & {
	children: LocationOption[];
};

export type CreateLocationRequest = {
	name: string;
	description?: string;
	parent_id?: number | null;
};

export type UpdateLocationRequest = {
	name?: string;
	description?: string;
	parent_id?: number | null;
};

export type ItemRequest = {
	user_id: number;
	item_id: number;
	created_at: string;
	status: ItemRequestStatus;
	reason: string;
	user_name?: string;
	item_name?: string;
};

export type ItemRequestsPage = {
	items: ItemRequestSummary[];
	limit: number;
	offset: number;
	total: number;
};

export type ItemRequestSummary = ItemRequest & {
	user_name: string;
	item_name: string;
};

export type ItemRequestUserOption = {
	id: number;
	name: string;
};

export type ItemRequestPreparationReport = {
	user: ItemRequestUserOption;
	items: ItemRequestPreparationItem[];
};

export type ItemRequestPreparationReportParams = Pick<
	ListItemRequestsParams,
	'status' | 'createdFrom' | 'createdTo'
> & {
	itemIDs?: number[];
};

export type ItemRequestPreparationItem = {
	id: number;
	name: string;
	type_name: string;
	location_names?: string[];
	derived_name_format: string;
	status: ItemRequestStatus;
	consumption: ConsumptionStatus;
	reason: string;
	requested_at: string;
	properties: ItemRequestPreparationProperty[];
};

export type ItemRequestPreparationProperty = {
	name: string;
	value: PropertyValue;
	value_type: PropertyValueType;
	visibility: PropertyVisibility;
	position: number;
};

export type CreatePersonalItemRequest = {
	item_id: number;
	reason: string;
};

export type ListItemRequestsParams = {
	limit?: number;
	offset?: number;
	status?: ItemRequestStatus;
	userIDs?: number[];
	createdFrom?: string;
	createdTo?: string;
};

export type NotificationKind = 'item_request' | 'item_expiry' | 'item_low_stock';
export type NotificationExpiryType = 'expiring_soon' | 'expired';

export type NotificationItemExpiry = {
	item: Item;
	expiry_type: NotificationExpiryType;
};

// Everything a low stock warning says was recorded when it fired: a stock group is a rendered name
// with no row behind it, and the type may since have been renamed or had its threshold moved.
export type NotificationLowStock = {
	type_id: number | null;
	type_name: string;
	group_name: string;
	threshold: number;
	observed: number;
};

// Only the descriptor matching `kind` is sent; all are absent if the descriptor row is missing.
export type Notification = {
	id: number;
	recipient_id: number;
	created_at: string;
	kind: NotificationKind;
	read: boolean;
	desc_item_request?: ItemRequest;
	desc_item_expiry?: NotificationItemExpiry;
	desc_low_stock?: NotificationLowStock;
};

export type NotificationsPage = {
	items: Notification[];
	limit: number;
	offset: number;
	total: number;
	total_unread: number;
};

export type ListNotificationsParams = {
	limit?: number;
	offset?: number;
};

export type AuditAction =
	| 'item_create'
	| 'item_update'
	| 'item_consume'
	| 'item_delete'
	| 'item_property_add'
	| 'item_property_update'
	| 'item_property_remove'
	| 'item_request_create'
	| 'item_request_approve'
	| 'item_request_delete'
	| 'item_type_create'
	| 'item_type_update'
	| 'item_type_delete'
	| 'item_type_property_add'
	| 'item_type_property_update'
	| 'item_type_property_remove'
	| 'item_type_property_reorder'
	| 'property_create'
	| 'property_update'
	| 'property_delete';

export type AuditTargetType = 'item' | 'item_type' | 'property';

// One field that moved. old and new stay raw JSON so a recorded property value can be rendered with
// the same helper a live one goes through. Nothing here is worded for a reader: key and value_type
// are identifiers the frontend turns into Serbian.
export type AuditChange = {
	key: string;
	// The recorded name of a user-defined property, when one entry spans several of them. Absent for
	// the fixed columns, which key already identifies.
	label?: string;
	value_type?: string;
	old?: unknown;
	new?: unknown;
};

export type AuditProperty = {
	id: number;
	name: string;
	value_type: string;
};

export type AuditContext = {
	type_id?: number;
	type_name?: string;
	batch_size?: number;
	property_id?: number;
	property_name?: string;
	subject_user_id?: number;
	subject_user_name?: string;
	reason?: string;
	request_status?: string;
	properties?: AuditProperty[];
};

export type AuditEntry = {
	id: number;
	created_at: string;
	// Null once that user is deleted, while actor_name still holds who it was. Both empty means there
	// was no user behind the change at all - the seeder, or a background job.
	actor_id: number | null;
	actor_name: string;
	action: AuditAction;
	target_type: AuditTargetType;
	// Still set after the target is gone, which is why an entry outlives what it describes.
	target_id: number;
	target_label: string;
	changes: AuditChange[];
	context: AuditContext;
};

export type AuditLogPage = {
	items: AuditEntry[];
	limit: number;
	offset: number;
	total: number;
};

export type AuditActorOption = {
	id: number;
	name: string;
};

export type ListAuditLogParams = {
	limit?: number;
	offset?: number;
	action?: AuditAction;
	targetType?: AuditTargetType;
	targetID?: number;
	actorID?: number;
	createdFrom?: string;
	createdTo?: string;
};

// The whole set, not a convenience list: anything else is rejected.
export type DashboardMonths = 3 | 6 | 12;

export type DashboardMonthCount = {
	// The first day of the bucket, YYYY-MM-DD, so it can be given straight to a date formatter.
	month: string;
	count: number;
};

// One month of consumption, already split by the state each item was left in.
export type DashboardConsumptionBucket = {
	month: string;
	fully_consumed: number;
	partially_consumed: number;
	damaged: number;
};

// A stock line names its own type, since a group name only means anything inside one.
export type DashboardStockGroup = {
	type_id: number;
	type_name: string;
	name: string;
	in_stock_count: number;
	total_count: number;
	low_stock_count: number | null;
	low: boolean;
};

// days_remaining is negative once the date has passed, which is the difference between a warning
// and a write-off.
export type DashboardExpiringItem = {
	id: number;
	name: string;
	type_name: string;
	expires_on: string;
	days_remaining: number;
};

// The two ranges read in opposite directions: expiry looks forward, consumption back. stock_groups
// is every group; the shortages are the ones marked `low`.
export type Dashboard = {
	months: DashboardMonths;
	type_id: number | null;
	expiring_by_month: DashboardMonthCount[];
	expired_backlog: number;
	// The list is capped; the total says how many there were before the cap.
	expiring_items: DashboardExpiringItem[];
	expiring_items_total: number;
	consumption_by_month: DashboardConsumptionBucket[];
	stock_groups: DashboardStockGroup[];
};

export type GetDashboardParams = {
	months?: DashboardMonths;
	typeID?: number;
};
