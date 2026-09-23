import type { Dashboard } from '$lib/api';

export const dashboardReportPrintContextKey = Symbol('dashboard-report-print');

export type DashboardReportPrintContext = {
	setDashboardReport: (report: Dashboard | null) => void;
	print: () => void;
};
