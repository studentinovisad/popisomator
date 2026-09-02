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
	role?: UserRole;
	status?: UserStatus;
};

export type ConsumptionStatus =
	'not_consumed' | 'partially_consumed' | 'fully_consumed' | 'damaged';

export type PropertyVisibility = 'overview' | 'details';

export type ItemProperty = {
	id: number;
	value: {};
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
};

export type ItemsPage = {
	items: Item[];
	limit: number;
	offset: number;
	total: number;
	// Sums of the structured properties across every item matching the filters, not just this page.
	totals: ItemPropertyTotal[];
};

// Sum of one structured property (price, mass, volume) across every item matching a filter. `value`
// has the same shape as the property itself - a PTPrice, or a PTMeasure already converted to a
// readable unit - so displayJson renders it like any other value. Prices are reported per currency,
// so one property can appear more than once.
export type ItemPropertyTotal = {
	property_id: number;
	value_type: PropertyValueType;
	value: {};
	// How many of the matched items carried the property; compare with ItemsPage.total.
	value_count: number;
};

export type ListItemsParams = PageRequest & {
	typeID?: number;
	propertyFilters?: Record<number, {}>;
};

export type UpdateItemRequest = {
	type_id?: number;
	consumption?: ConsumptionStatus;
};

export type CreateItemRequest = {
	type_id: number;
	properties: ItemProperty[];
	amount: number;
};

export type ItemTypeProperty = {
	id: number;
	default_value: {} | null;
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
};

export type UpdateItemTypeRequest = Partial<
	Pick<ItemType, 'name' | 'description' | 'derived_name_format'>
>;

export type AddUpdateItemTypePropertyRequest = {
	property_id?: number;
	default_value?: {};
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
	default_value: {} | null;
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
	default_value: {} | null;
};

export type UpdatePropertyRequest = Partial<
	Pick<CreatePropertyRequest, 'name' | 'description' | 'default_value'>
>;

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

export type CreatePersonalItemRequest = {
	item_id: number;
	reason: string;
};

export type ListItemRequestsParams = {
	limit?: number;
	offset?: number;
	status?: ItemRequestStatus;
};
