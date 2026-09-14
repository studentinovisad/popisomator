import type { ItemRequestPreparationReport } from '$lib/api';

export const preparationReportPrintContextKey = Symbol('preparation-report-print');

export type PreparationReportPrintContext = {
	setPreparationReports: (reports: ItemRequestPreparationReport[]) => void;
	print: () => void;
};
