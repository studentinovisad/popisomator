import { request } from '$lib/api/client';
import type { Dashboard, GetDashboardParams } from '$lib/api/types';

export const dashboardApi = {
	getDashboard: ({ months, typeID }: GetDashboardParams = {}) => {
		const query = new URLSearchParams();
		if (months) query.set('months', String(months));
		if (typeID) query.set('type_id', String(typeID));
		return request<Dashboard>(`/dashboard?${query}`);
	}
};
