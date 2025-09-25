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
  private readonly apiUrl = `${environment.accountServiceUrl}/accounts`;

  // CRUD Operations
  createCard(cardData: CreateCardRequest): Observable<Card> {
    // Encriptar datos sensibles antes de enviar
    return from(this.encryptCardData(cardData)).pipe(
      switchMap(encryptedData => {
        // Usar el endpoint de cards con el accountId del usuario
        const cardPayload = {
          card_type: cardData.cardType,
          card_brand: this.detectCardBrand(cardData.cardNumber),
          last_four_digits: this.getLastFourDigits(cardData.cardNumber),
          masked_number: this.maskCardNumber(cardData.cardNumber),
          holder_name: cardData.holderName,
          expiration_month: cardData.expirationMonth,
          expiration_year: cardData.expirationYear,
          nickname: cardData.nickname,
          encrypted_number: encryptedData.encryptedNumber,
          encrypted_cvv: encryptedData.encryptedCvv,
          key_fingerprint: encryptedData.keyFingerprint,
          // Campos específicos para tarjetas de crédito
          ...(cardData.cardType === CardType.CREDIT && cardData.creditLimit && {
            credit_limit: cardData.creditLimit
          }),
          ...(cardData.closingDate && {
            closing_date: new Date(cardData.closingDate).toISOString().split('T')[0]
          }),
          ...(cardData.dueDate && {
            due_date: new Date(cardData.dueDate).toISOString().split('T')[0]
          })
        };

        return this.http.post<any>(`${this.apiUrl}/${cardData.accountId}/cards`, cardPayload).pipe(
          map(response => this.mapCardResponseToCard(response)),
          catchError(error => {
            console.error('Error creating card:', error);
            throw new Error('Error al crear la tarjeta: ' + (error.error?.message || error.message));
          })
        );
      })
    );
  }

  getCardsByAccount(accountId: string, page: number = 1, pageSize: number = 20): Observable<CardsListResponse> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());
    
    return this.http.get<any>(`${this.apiUrl}/${accountId}/cards`, { params }).pipe(
      map(response => {
        const cards = response.data ? response.data.map((card: any) => this.mapCardResponseToCard(card)) : [];
        
        return {
          cards,
          total: response.pagination?.total_items || cards.length,
          page: response.pagination?.current_page || page,
          pageSize: response.pagination?.page_size || pageSize,
          totalPages: response.pagination?.total_pages || Math.ceil(cards.length / pageSize)
        };
      }),
      catchError(error => {
        console.error('Error getting cards by account:', error);
        return of({
          cards: [],
          total: 0,
          page: page,
          pageSize: pageSize,
          totalPages: 0
        });
      })
    );
  }

  getCardsByUser(userId: string, page: number = 1, pageSize: number = 20): Observable<CardsListResponse> {
    console.log('CardService: Getting cards for user', userId);
    
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());
    
    return this.http.get<any>(`${this.apiUrl}/user/${userId}/cards`, { params }).pipe(
      map(response => {
        console.log('CardService: Received response from backend:', response);
        
        const cards = response.data ? response.data.map((card: any) => this.mapCardResponseToCard(card)) : [];
        
        console.log('CardService: Mapped cards:', cards);
        
        return {
          cards,
          total: response.pagination?.total_items || cards.length,
          page: response.pagination?.current_page || page,
          pageSize: response.pagination?.page_size || pageSize,
          totalPages: response.pagination?.total_pages || Math.ceil(cards.length / pageSize)
        };
      }),
      catchError(error => {
        console.error('CardService: Error getting cards by user:', error);
        return of({
          cards: [],
          total: 0,
          page: page,
          pageSize: pageSize,
          totalPages: 0
        });
      })
    );
  }

  getCardById(accountId: string, cardId: string): Observable<Card> {
    return this.http.get<any>(`${this.apiUrl}/${accountId}/cards/${cardId}`).pipe(
      map(response => this.mapCardResponseToCard(response))
    );
  }

  updateCard(accountId: string, cardId: string, cardData: UpdateCardRequest): Observable<Card> {
    const updatePayload = {
      holder_name: cardData.holderName,
      nickname: cardData.nickname,
      expiration_month: cardData.expirationMonth,
      expiration_year: cardData.expirationYear
    };

    return this.http.put<any>(`${this.apiUrl}/${accountId}/cards/${cardId}`, updatePayload).pipe(
      map(response => this.mapCardResponseToCard(response))
    );
  }

  deleteCard(accountId: string, cardId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${accountId}/cards/${cardId}`);
  }

  setDefaultCard(accountId: string, cardId: string): Observable<Card> {
    return this.http.put<any>(`${this.apiUrl}/${accountId}/cards/${cardId}/set-default`, {}).pipe(
      map(response => this.mapCardResponseToCard(response))
    );
  }

  blockCard(accountId: string, cardId: string): Observable<Card> {
    return this.http.put<any>(`${this.apiUrl}/${accountId}/cards/${cardId}/block`, {}).pipe(
      map(response => this.mapCardResponseToCard(response))
    );
  }

  unblockCard(accountId: string, cardId: string): Observable<Card> {
    return this.http.put<any>(`${this.apiUrl}/${accountId}/cards/${cardId}/unblock`, {}).pipe(
      map(response => this.mapCardResponseToCard(response))
    );
  }

  // Métodos específicos para activar/desactivar tarjetas (usar los nuevos endpoints)
  activateCard(accountId: string, cardId: string): Observable<Card> {
    return this.unblockCard(accountId, cardId);
  }

  deactivateCard(accountId: string, cardId: string): Observable<Card> {
    return this.blockCard(accountId, cardId);
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
    
    // Usar siempre el formato estándar: **** **** **** XXXX (19 caracteres)
    // Esto evita problemas con números de tarjeta de diferentes longitudes
    return `**** **** **** ${lastFour}`;
  }

  getLastFourDigits(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    return cleanNumber.slice(-4);
  }

  formatCardNumber(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    return cleanNumber.match(/.{1,4}/g)?.join(' ') || cardNumber;
  }

  // Método para mapear la respuesta del backend (Card API) a Card del frontend
  private mapCardResponseToCard(cardResponse: any): Card {
    return {
      id: cardResponse.id,
      accountId: cardResponse.account_id,
      cardType: cardResponse.card_type as CardType,
      cardBrand: cardResponse.card_brand as CardBrand,
      lastFourDigits: cardResponse.last_four_digits || cardResponse.lastFourDigits,
      maskedNumber: cardResponse.masked_number || cardResponse.maskedNumber,
      holderName: cardResponse.holder_name || cardResponse.holderName,
      expirationMonth: cardResponse.expiration_month || cardResponse.expirationMonth,
      expirationYear: cardResponse.expiration_year || cardResponse.expirationYear,
      status: cardResponse.status as CardStatus,
      isDefault: cardResponse.is_default || cardResponse.isDefault,
      nickname: cardResponse.nickname,
      balance: cardResponse.balance || 0, // New: Include balance field
      creditLimit: cardResponse.credit_limit,
      closingDate: cardResponse.closing_date,
      dueDate: cardResponse.due_date,
      createdAt: cardResponse.created_at,
      updatedAt: cardResponse.updated_at
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