import { catalogApi } from '$lib/api/catalog';
import { itemRequestsApi } from '$lib/api/item-requests';
import { itemsApi } from '$lib/api/items';
import { usersApi } from '$lib/api/users';

export { ApiError } from '$lib/api/client';
export type * from '$lib/api/types';

export const api = {
	...usersApi,
	...itemsApi,
	...catalogApi,
	...itemRequestsApi
};
