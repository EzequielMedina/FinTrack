import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface ExchangeRate {
  compra: number;
  venta: number;
  casa: string;
  nombre: string;
  moneda: string;
  fechaActualizacion: Date;
  spread: number;
  spreadPercentage: number;
}

export interface ExchangeRateResponse {
  compra: number;
  venta: number;
  casa: string;
  nombre: string;
  moneda: string;
  fechaActualizacion: string;
  spread: number;
  spreadPercentage: number;
}

@Injectable({
  providedIn: 'root'
})
export class ExchangeService {
  private readonly apiUrl = `${environment.apiUrl}/exchange`;

  constructor(private http: HttpClient) {}

  /**
   * Obtiene la cotización del dólar oficial
   */
  getDolarOficial(): Observable<ExchangeRate> {
    return this.http.get<ExchangeRate>(`${this.apiUrl}/dolar-oficial`);
  }

  /**
   * Convierte un monto de ARS a USD usando la cotización actual
   */
  convertArsToUsd(amount: number, rate: ExchangeRate): number {
    return amount / rate.venta;
  }

  /**
   * Convierte un monto de USD a ARS usando la cotización actual
   */
  convertUsdToArs(amount: number, rate: ExchangeRate): number {
    return amount * rate.compra;
  }

  /**
   * Formatea la cotización para mostrar
   */
  formatExchangeRate(rate: ExchangeRate): string {
    return `Compra: $${rate.compra.toLocaleString('es-AR', { minimumFractionDigits: 2 })} - Venta: $${rate.venta.toLocaleString('es-AR', { minimumFractionDigits: 2 })}`;
  }
}