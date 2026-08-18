export type PageRequest = {
	limit?: number;
	offset?: number;
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
	value: string;
	visibility?: PropertyVisibility;
};

export type Item = {
	id: number;
	consumption: ConsumptionStatus;
	properties: ItemProperty[];
	type_id: number;
};

export type ItemsPage = {
	items: Item[];
	limit: number;
	offset: number;
	total: number;
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
	default_value: string | null;
	name?: string;
	visibility: PropertyVisibility;
};

export type ItemTypeOption = {
	id: number;
	name: string;
	description: string;
	properties: ItemTypeProperty[];
};

export type ItemType = ItemTypeOption;

export type ItemTypesPage = {
	items: ItemType[];
	limit: number;
	offset: number;
	total: number;
};

export type CreateItemTypeRequest = {
	name: string;
	description: string;
	properties: ItemTypeProperty[];
};

export type UpdateItemTypeRequest = Partial<Pick<ItemType, 'name' | 'description'>>;

export type AddUpdateItemTypePropertyRequest = {
	property_id?: number;
	default_value?: string;
	visibility?: string;
};

export type PropertyValueType = 'string' | 'number' | 'boolean';

export type PropertyOption = {
	id: number;
	name: string;
	value_type: PropertyValueType;
	description: string;
	default_value: string | null;
};

export type Property = PropertyOption;

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
	default_value: string | null;
};

export type UpdatePropertyRequest = Partial<
	Pick<CreatePropertyRequest, 'name' | 'description' | 'default_value'>
>;
