export enum InstallmentStatus {
  PENDING = 'pending',
  PAID = 'paid',
  OVERDUE = 'overdue',
  CANCELLED = 'cancelled'
}

export enum InstallmentPlanStatus {
  ACTIVE = 'active',
  COMPLETED = 'completed',
  CANCELLED = 'cancelled',
  SUSPENDED = 'suspended'
}

export interface Installment {
  id: string;
  planId: string;
  installmentNumber: number;
  amount: number;
  dueDate: string;
  status: InstallmentStatus;
  paidDate?: string;
  paidAmount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface InstallmentPlan {
  id: string;
  transactionId: string;
  cardId: string;
  userId: string;
  
  // Plan details
  totalAmount: number;
  installmentsCount: number;
  installmentAmount: number;
  startDate: string;
  
  // Merchant information
  merchantName?: string;
  merchantId?: string;
  description?: string;
  
  // Status tracking
  status: InstallmentPlanStatus;
  paidInstallments: number;
  remainingAmount: number;
  
  // Interest and fees
  interestRate: number;
  totalInterest: number;
  adminFee: number;
  
  // Audit fields
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
  cancelledAt?: string;
  
  // Relationships
  card?: any; // Will be typed as Card when needed
  installments?: Installment[];
}

export interface InstallmentPreview {
  totalAmount: number;
  installmentsCount: number;
  installmentAmount: number;
  startDate: string;
  interestRate: number;
  totalInterest: number;
  adminFee: number;
  totalToPay: number;
  installments: InstallmentPreviewItem[];
}

export interface InstallmentPreviewItem {
  number: number;
  amount: number;
  dueDate: string;
  principal: number;
  interest: number;
  fee: number;
}

export interface CreateInstallmentPlanRequest {
  cardId: string;
  totalAmount: number;
  installmentsCount: number;
  startDate: string;
  description?: string;
  merchantName?: string;
  merchantId?: string;
  interestRate?: number;
  adminFee?: number;
  reference?: string;
}

export interface InstallmentPreviewRequest {
  amount: number;
  installmentsCount: number;
  startDate: string;
  interestRate?: number;
  adminFee?: number;
}

export interface PayInstallmentRequest {
  installmentId: string;
  amount: number;
  paymentMethod?: string;
  reference?: string;
}

export interface ChargeWithInstallmentsRequest {
  cardId: string;
  amount: number;
  totalAmount: number;
  installmentsCount: number;
  startDate: string;
  description: string;
  merchantName?: string;
  reference?: string;
  interestRate?: number;
  adminFee?: number;
}

export interface ChargeWithInstallmentsResponse {
  installmentPlan: InstallmentPlan;
  card: any; // Will be typed as Card when needed
  firstInstallmentCharged: boolean;
  transactionId: string;
}

export interface InstallmentPlanResponse {
  installmentPlan: InstallmentPlan;
  installments: Installment[];
  upcomingPayment?: Installment;
  overdueCount: number;
  totalPaid: number;
  totalRemaining: number;
}

export interface InstallmentSummary {
  userId: string;
  totalPlans: number;
  activePlans: number;
  completedPlans: number;
  cancelledPlans: number;
  totalInstallments: number;
  pendingInstallments: number;
  paidInstallments: number;
  overdueInstallments: number;
  totalAmountOwed: number;
  totalAmountPaid: number;
  nextPaymentDue?: string;
  nextPaymentAmount: number;
}

export interface MonthlyInstallmentLoad {
  year: number;
  month: number;
  totalInstallments: number;
  totalAmount: number;
  paidAmount: number;
  pendingAmount: number;
}

export interface InstallmentPlansListResponse {
  plans: InstallmentPlan[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface InstallmentsListResponse {
  installments: Installment[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Error types
export interface InstallmentError {
  code: string;
  message: string;
  field?: string;
}

export interface InstallmentValidationError {
  field: string;
  message: string;
}

export interface InstallmentValidationResult {
  isValid: boolean;
  errors: InstallmentValidationError[];
}

// Form types
export interface InstallmentCalculatorForm {
  amount: number;
  installmentsCount: number;
  startDate: string;
  cardId?: string;
}

export interface InstallmentPaymentForm {
  installmentId: string;
  amount: number;
  paymentMethod: string;
  notes?: string;
}

export interface InstallmentFormErrors {
  amount?: string;
  installmentsCount?: string;
  startDate?: string;
  cardId?: string;
  paymentMethod?: string;
}

// Configuration types
export interface InstallmentConfig {
  maxInstallments: number;
  minAmount: number;
  maxAmount: number;
  defaultInterestRate: number;
  adminFeePercentage: number;
  allowedInstallmentCounts: number[];
}

// Chart/Display types
export interface InstallmentChartData {
  labels: string[];
  principal: number[];
  interest: number[];
  fees: number[];
}

export interface InstallmentCalendarItem {
  date: string;
  installments: Installment[];
  totalAmount: number;
  isOverdue: boolean;
}