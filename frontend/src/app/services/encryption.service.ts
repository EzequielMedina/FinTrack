import { Injectable } from '@angular/core';

export interface EncryptionResult {
  encryptedData: string;
  keyFingerprint: string;
  algorithm: string;
}

@Injectable({
  providedIn: 'root'
})
export class EncryptionService {
  private readonly algorithm = 'AES-GCM';
  private readonly keyLength = 256;
  
  // En producción, estas claves deberían venir del servidor
  private readonly publicKey = 'dev-public-key-base64';
  private readonly keyFingerprint = 'dev-key-fingerprint-' + Date.now();

  /**
   * Encripta datos sensibles usando Web Crypto API
   * IMPORTANTE: Esta es una implementación básica para desarrollo
   * En producción, usar claves del servidor y protocolos seguros
   */
  async encryptSensitiveData(data: string): Promise<EncryptionResult> {
    try {
      // En desarrollo, usamos base64 encoding (NO SEGURO PARA PRODUCCIÓN)
      if (this.isDevelopment()) {
        return this.mockEncryption(data);
      }

      // Implementación real para producción
      return await this.realEncryption(data);
    } catch (error) {
      console.error('Encryption error:', error);
      throw new Error('Error al encriptar datos sensibles');
    }
  }

  /**
   * Valida el formato de un número de tarjeta
   */
  validateCardNumberFormat(cardNumber: string): boolean {
    // Remover espacios y guiones
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    
    // Validar que solo contenga números
    if (!/^\d+$/.test(cleanNumber)) {
      return false;
    }

    // Validar longitud (13-19 dígitos)
    if (cleanNumber.length < 13 || cleanNumber.length > 19) {
      return false;
    }

    return true;
  }

  /**
   * Algoritmo de Luhn para validar números de tarjeta
   */
  validateLuhnAlgorithm(cardNumber: string): boolean {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    let sum = 0;
    let shouldDouble = false;

    for (let i = cleanNumber.length - 1; i >= 0; i--) {
      let digit = parseInt(cleanNumber.charAt(i));

      if (shouldDouble) {
        digit *= 2;
        if (digit > 9) {
          digit -= 9;
        }
      }

      sum += digit;
      shouldDouble = !shouldDouble;
    }

    return sum % 10 === 0;
  }

  /**
   * Detecta el tipo de tarjeta basado en el número
   */
  detectCardBrand(cardNumber: string): string {
    const cleanNumber = cardNumber.replace(/[\s-]/g, '');
    
    // Patrones de detección de marca
    const patterns = {
      visa: /^4[0-9]{12}(?:[0-9]{3})?$/,
      mastercard: /^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$/,
      amex: /^3[47][0-9]{13}$/,
      discover: /^6(?:011|5[0-9]{2})[0-9]{12}$/,
      diners: /^3[0689][0-9]{13}$/
    };

    for (const [brand, pattern] of Object.entries(patterns)) {
      if (pattern.test(cleanNumber)) {
        return brand;
      }
    }

    return 'unknown';
  }

  /**
   * Valida CVV según el tipo de tarjeta
   */
  validateCVV(cvv: string, cardBrand: string): boolean {
    if (!/^\d+$/.test(cvv)) {
      return false;
    }

    // American Express usa 4 dígitos, otros usan 3
    const expectedLength = cardBrand === 'amex' ? 4 : 3;
    return cvv.length === expectedLength;
  }

  /**
   * Genera un hash seguro para almacenamiento local temporal
   */
  async generateSecureHash(data: string): Promise<string> {
    const encoder = new TextEncoder();
    const dataBuffer = encoder.encode(data);
    
    if (crypto.subtle) {
      const hashBuffer = await crypto.subtle.digest('SHA-256', dataBuffer);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    }

    // Fallback para entornos sin Web Crypto API
    return btoa(data).slice(0, 32);
  }

  /**
   * Limpia datos sensibles de la memoria
   */
  clearSensitiveData(data: any): void {
    if (typeof data === 'object' && data !== null) {
      Object.keys(data).forEach(key => {
        if (typeof data[key] === 'string') {
          data[key] = '';
        } else if (typeof data[key] === 'object') {
          this.clearSensitiveData(data[key]);
        }
      });
    }
  }

  private isDevelopment(): boolean {
    return !window.location.hostname.includes('production-domain.com');
  }

  private mockEncryption(data: string): EncryptionResult {
    console.warn('⚠️  DESARROLLO: Usando encriptación simulada. NO usar en producción.');
    
    return {
      encryptedData: btoa(data), // Base64 - NO SEGURO
      keyFingerprint: this.keyFingerprint,
      algorithm: 'MOCK-BASE64'
    };
  }

  private async realEncryption(data: string): Promise<EncryptionResult> {
    // Implementación real usando Web Crypto API
    // En producción, la clave pública vendría del servidor
    
    const encoder = new TextEncoder();
    const dataBuffer = encoder.encode(data);
    
    // Generar clave temporal (en producción usar clave del servidor)
    const key = await crypto.subtle.generateKey(
      {
        name: this.algorithm,
        length: this.keyLength
      },
      false, // no exportable
      ['encrypt']
    );

    // Generar IV único
    const iv = crypto.getRandomValues(new Uint8Array(12));
    
    // Encriptar
    const encryptedBuffer = await crypto.subtle.encrypt(
      {
        name: this.algorithm,
        iv: iv
      },
      key,
      dataBuffer
    );

    // Combinar IV + datos encriptados
    const combinedBuffer = new Uint8Array(iv.length + encryptedBuffer.byteLength);
    combinedBuffer.set(iv);
    combinedBuffer.set(new Uint8Array(encryptedBuffer), iv.length);

    // Convertir a base64
    const encryptedData = btoa(String.fromCharCode(...combinedBuffer));

    return {
      encryptedData,
      keyFingerprint: this.keyFingerprint,
      algorithm: this.algorithm
    };
  }
}