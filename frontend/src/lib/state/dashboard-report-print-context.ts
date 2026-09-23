import type { Dashboard } from '$lib/api';

export const dashboardReportPrintContextKey = Symbol('dashboard-report-print');

export type DashboardReportPrintMode = 'trends' | 'stock';

export type DashboardReportPrintContext = {
	setDashboardReport: (report: Dashboard | null) => void;
	setPrintMode: (mode: DashboardReportPrintMode | null) => void;
	print: () => void;
};
