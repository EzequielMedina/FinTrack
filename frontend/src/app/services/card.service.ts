import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, of, from } from 'rxjs';
import { map, catchError, switchMap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { EncryptionService } from './encryption.service';
import {
  Card,
  CreateCardRequest,
  UpdateCardRequest,
  CardsListResponse,
  CardValidationResult,
  CardValidationError,
  CardBrand,
  CardType,
  CardStatus,
  EncryptedCardData
} from '../models';

@Injectable({ providedIn: 'root' })
export class CardService {
  private readonly http = inject(HttpClient);
  private readonly encryptionService = inject(EncryptionService);
  private readonly apiUrl = environment.accountServiceUrl;

  // CRUD Operations
  createCard(cardData: CreateCardRequest): Observable<Card> {
    // Encriptar datos sensibles antes de enviar
    return from(this.encryptCardData(cardData)).pipe(
      switchMap(encryptedData => {
        // Mapear Card a Account para el backend
        const accountPayload = {
          user_id: cardData.userId,
          account_type: cardData.cardType, // 'credit' o 'debit'
          name: cardData.nickname || `Tarjeta ${cardData.cardType === CardType.CREDIT ? 'de Crédito' : 'de Débito'}`,
          description: `${cardData.holderName} - **** **** **** ${cardData.cardNumber.slice(-4)} - Exp: ${cardData.expirationMonth.toString().padStart(2, '0')}/${cardData.expirationYear}`,
          currency: cardData.currency || 'ARS',
          initial_balance: 0.0
        };

                return this.http.post<any>(this.apiUrl, accountPayload).pipe(
          map(accountResponse => this.mapAccountToCard(accountResponse, cardData))
        );
      })
    );
  }

  getCardsByAccount(accountId: string, page: number = 1, pageSize: number = 20): Observable<CardsListResponse> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());
    
    // Obtener todas las cuentas y filtrar por tipo tarjeta
    return this.http.get<any>(`${this.apiUrl}`, { params }).pipe(
      map(response => {
        // Filtrar solo cuentas de tipo credit y debit
        const cardAccounts = response.data.filter((account: any) => 
          account.account_type === 'credit' || account.account_type === 'debit'
        );
        
        const cards = cardAccounts.map((account: any) => this.mapAccountToCard(account));
        
        return {
          cards,
          total: cardAccounts.length,
          page: response.pagination.current_page,
          pageSize: response.pagination.page_size,
          totalPages: Math.ceil(cardAccounts.length / pageSize)
        };
      })
    );
  }

  getCardsByUser(userId: string, page: number = 1, pageSize: number = 20): Observable<CardsListResponse> {
    console.log('CardService: Getting cards for user', userId);
    
    // Obtener cuentas por usuario y filtrar tarjetas desde account-service
    return this.http.get<any[]>(`${this.apiUrl}/user/${userId}`).pipe(
      map(accounts => {
        console.log('CardService: Received accounts from backend:', accounts);
        
        // Filtrar solo cuentas de tipo credit y debit
        const cardAccounts = accounts.filter((account: any) => 
          account.account_type === 'credit' || account.account_type === 'debit'
        ) || [];
        
        console.log('CardService: Filtered card accounts:', cardAccounts);
        
        const cards = cardAccounts.map((account: any) => this.mapAccountToCard(account));
        console.log('CardService: Mapped cards:', cards);
        
        return {
          cards,
          total: cardAccounts.length,
          page: page,
          pageSize: pageSize,
          totalPages: Math.ceil(cardAccounts.length / pageSize)
        };
      })
    );
  }

  getCardById(accountId: string, cardId: string): Observable<Card> {
    return this.http.get<any>(`${this.apiUrl}/${cardId}`).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  updateCard(accountId: string, cardId: string, cardData: UpdateCardRequest): Observable<Card> {
    // Mapear UpdateCardRequest a UpdateAccountRequest
    const updatePayload = {
      name: cardData.nickname || `Tarjeta ${cardData.holderName ? 'de ' + cardData.holderName : ''}`,
      description: cardData.holderName || ''
    };

    return this.http.put<any>(`${this.apiUrl}/${cardId}`, updatePayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  deleteCard(accountId: string, cardId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${cardId}`);
  }

  setDefaultCard(accountId: string, cardId: string): Observable<Card> {
    // Por ahora, simulamos el comportamiento actualizando el estado
    const updatePayload = {
      name: "Tarjeta Predeterminada",
      description: "Tarjeta establecida como predeterminada"
    };
    
    return this.http.put<any>(`${this.apiUrl}/${cardId}`, updatePayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  blockCard(accountId: string, cardId: string): Observable<Card> {
    // Usar el endpoint de status del account-service
    const statusPayload = {
      is_active: false
    };
    
    return this.http.put<any>(`${this.apiUrl}/${cardId}/status`, statusPayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  unblockCard(accountId: string, cardId: string): Observable<Card> {
    // Usar el endpoint de status del account-service
    const statusPayload = {
      is_active: true
    };
    
    return this.http.put<any>(`${this.apiUrl}/${cardId}/status`, statusPayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  // Métodos específicos para activar/desactivar tarjetas
  activateCard(cardId: string): Observable<Card> {
    const statusPayload = {
      is_active: true
    };
    
    return this.http.put<any>(`${this.apiUrl}/${cardId}/status`, statusPayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  deactivateCard(cardId: string): Observable<Card> {
    const statusPayload = {
      is_active: false
    };
    
    return this.http.put<any>(`${this.apiUrl}/${cardId}/status`, statusPayload).pipe(
      map(accountResponse => this.mapAccountToCard(accountResponse))
    );
  }

  // Validation Methods
  validateCardNumber(cardNumber: string): CardValidationResult {
    const errors: CardValidationError[] = [];
    let brand: CardBrand | undefined;

    // Remover espacios y guiones
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');

    // Validar formato básico
    if (!this.encryptionService.validateCardNumberFormat(cleanNumber)) {
      errors.push({
        field: 'cardNumber',
        message: 'El número de tarjeta debe contener entre 13 y 19 dígitos'
      });
    }

    // Detectar marca de tarjeta
    const detectedBrand = this.encryptionService.detectCardBrand(cleanNumber);
    brand = this.mapBrandStringToEnum(detectedBrand);

    // Validar longitud según la marca
    if (!this.isValidCardLength(cleanNumber, brand)) {
      errors.push({
        field: 'cardNumber',
        message: 'Número de tarjeta inválido para la marca detectada'
      });
    }

    // Algoritmo de Luhn
    if (!this.encryptionService.validateLuhnAlgorithm(cleanNumber)) {
      errors.push({
        field: 'cardNumber',
        message: 'El número de tarjeta no es válido'
      });
    }

    return {
      isValid: errors.length === 0,
      errors,
      brand
    };
  }

  validateCVV(cvv: string, cardBrand?: CardBrand): CardValidationResult {
    const errors: CardValidationError[] = [];

    const brandString = cardBrand ? this.mapBrandEnumToString(cardBrand) : 'unknown';
    
    if (!this.encryptionService.validateCVV(cvv, brandString)) {
      const expectedLength = cardBrand === CardBrand.AMERICAN_EXPRESS ? 4 : 3;
      errors.push({
        field: 'cvv',
        message: `El CVV debe tener ${expectedLength} dígitos numéricos`
      });
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  validateExpirationDate(month: number, year: number): CardValidationResult {
    const errors: CardValidationError[] = [];
    const now = new Date();
    const currentYear = now.getFullYear();
    const currentMonth = now.getMonth() + 1;

    if (month < 1 || month > 12) {
      errors.push({
        field: 'expirationMonth',
        message: 'El mes debe estar entre 1 y 12'
      });
    }

    if (year < currentYear || (year === currentYear && month < currentMonth)) {
      errors.push({
        field: 'expirationYear',
        message: 'La fecha de expiración no puede ser en el pasado'
      });
    }

    if (year > currentYear + 20) {
      errors.push({
        field: 'expirationYear',
        message: 'La fecha de expiración es demasiado lejana'
      });
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  // Utility Methods
  detectCardBrand(cardNumber: string): CardBrand {
    const patterns = {
      [CardBrand.VISA]: /^4[0-9]{12}(?:[0-9]{3})?$/,
      [CardBrand.MASTERCARD]: /^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$/,
      [CardBrand.AMERICAN_EXPRESS]: /^3[47][0-9]{13}$/,
      [CardBrand.DISCOVER]: /^6(?:011|5[0-9]{2})[0-9]{12}$/,
      [CardBrand.DINERS]: /^3[0689][0-9]{13}$/
    };

    for (const [brand, pattern] of Object.entries(patterns)) {
      if (pattern.test(cardNumber)) {
        return brand as CardBrand;
      }
    }

    return CardBrand.OTHER;
  }

  maskCardNumber(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    if (cleanNumber.length < 4) return cardNumber;
    
    const lastFour = cleanNumber.slice(-4);
    const masked = '*'.repeat(cleanNumber.length - 4);
    
    // Formatear con espacios cada 4 dígitos
    const formatted = (masked + lastFour).match(/.{1,4}/g)?.join(' ') || '';
    return formatted;
  }

  getLastFourDigits(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    return cleanNumber.slice(-4);
  }

  formatCardNumber(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    return cleanNumber.match(/.{1,4}/g)?.join(' ') || cardNumber;
  }

  // Método privado para mapear Account (backend) a Card (frontend)
  private mapAccountToCard(account: any, originalCardData?: CreateCardRequest): Card {
    // Extraer los últimos 4 dígitos del description si está disponible
    const lastFourMatch = account.description?.match(/\*{4}\s\*{4}\s\*{4}\s(\d{4})/) || 
                         account.description?.match(/(\d{4})/) ||
                         account.name?.match(/(\d{4})$/);
    const lastFourDigits = lastFourMatch ? lastFourMatch[1] : '0000';
    
    // Extraer datos de expiración del description si están disponibles
    const expMatch = account.description?.match(/Exp:\s(\d{2})\/(\d{4})/);
    const expirationMonth = expMatch ? parseInt(expMatch[1]) : (originalCardData?.expirationMonth || 12);
    const expirationYear = expMatch ? parseInt(expMatch[2]) : (originalCardData?.expirationYear || 2030);
    
    // Extraer nombre del titular del description
    const holderNameMatch = account.description?.match(/^([^-]+)/);
    const holderName = holderNameMatch ? holderNameMatch[1].trim() : (originalCardData?.holderName || 'Titular');
    
    // Determinar el brand basado en el primer dígito (simplificado)
    const firstDigit = lastFourDigits[0];
    let cardBrand = CardBrand.OTHER;
    if (firstDigit === '4') cardBrand = CardBrand.VISA;
    else if (firstDigit === '5') cardBrand = CardBrand.MASTERCARD;
    else if (firstDigit === '3') cardBrand = CardBrand.AMERICAN_EXPRESS;
    
    return {
      id: account.id,
      accountId: account.id, // Usar el ID de la cuenta como accountId
      userId: account.user_id,
      cardType: account.account_type === 'credit' ? CardType.CREDIT : CardType.DEBIT,
      cardBrand: cardBrand,
      lastFourDigits: lastFourDigits,
      maskedNumber: `**** **** **** ${lastFourDigits}`,
      holderName: holderName,
      expirationMonth: expirationMonth,
      expirationYear: expirationYear,
      status: account.is_active ? CardStatus.ACTIVE : CardStatus.INACTIVE,
      isDefault: false, // Por ahora no manejamos tarjeta predeterminada
      nickname: account.name,
      currency: account.currency || 'ARS',
      balance: account.balance || 0,
      createdAt: account.created_at,
      updatedAt: account.updated_at
    };
  }

  getCardBrandIcon(brand: CardBrand): string {
    const icons = {
      [CardBrand.VISA]: 'assets/images/cards/visa.svg',
      [CardBrand.MASTERCARD]: 'assets/images/cards/mastercard.svg',
      [CardBrand.AMERICAN_EXPRESS]: 'assets/images/cards/amex.svg',
      [CardBrand.DISCOVER]: 'assets/images/cards/discover.svg',
      [CardBrand.DINERS]: 'assets/images/cards/diners.svg',
      [CardBrand.OTHER]: 'assets/images/cards/generic.svg'
    };

    return icons[brand];
  }

  // Private Methods
  private async encryptCardData(cardData: CreateCardRequest): Promise<EncryptedCardData> {
    try {
      const [encryptedNumber, encryptedCvv] = await Promise.all([
        this.encryptionService.encryptSensitiveData(cardData.cardNumber),
        this.encryptionService.encryptSensitiveData(cardData.cvv)
      ]);

      return {
        encryptedNumber: encryptedNumber.encryptedData,
        encryptedCvv: encryptedCvv.encryptedData,
        keyFingerprint: encryptedNumber.keyFingerprint
      };
    } catch (error) {
      console.error('Error encrypting card data:', error);
      throw new Error('Error al procesar datos de la tarjeta');
    }
  }

  private mapBrandStringToEnum(brandString: string): CardBrand {
    const mapping: { [key: string]: CardBrand } = {
      'visa': CardBrand.VISA,
      'mastercard': CardBrand.MASTERCARD,
      'amex': CardBrand.AMERICAN_EXPRESS,
      'discover': CardBrand.DISCOVER,
      'diners': CardBrand.DINERS,
      'unknown': CardBrand.OTHER
    };
    
    return mapping[brandString] || CardBrand.OTHER;
  }

  private mapBrandEnumToString(brand: CardBrand): string {
    const mapping: { [key in CardBrand]: string } = {
      [CardBrand.VISA]: 'visa',
      [CardBrand.MASTERCARD]: 'mastercard',
      [CardBrand.AMERICAN_EXPRESS]: 'amex',
      [CardBrand.DISCOVER]: 'discover',
      [CardBrand.DINERS]: 'diners',
      [CardBrand.OTHER]: 'unknown'
    };
    
    return mapping[brand];
  }

  private isValidCardLength(cardNumber: string, brand: CardBrand): boolean {
    const validLengths = {
      [CardBrand.VISA]: [13, 16, 19],
      [CardBrand.MASTERCARD]: [16],
      [CardBrand.AMERICAN_EXPRESS]: [15],
      [CardBrand.DISCOVER]: [16],
      [CardBrand.DINERS]: [14],
      [CardBrand.OTHER]: [13, 14, 15, 16, 17, 18, 19]
    };

    return validLengths[brand].includes(cardNumber.length);
  }
}