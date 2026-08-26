import type { ItemRequestPreparationReport } from '$lib/api';

export const preparationReportPrintContextKey = Symbol('preparation-report-print');

export type PreparationReportPrintContext = {
	setPreparationReport: (report: ItemRequestPreparationReport | null) => void;
	print: () => void;
};
