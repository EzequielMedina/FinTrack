import { Pipe, PipeTransform } from '@angular/core';

/**
 * Formatea un string de fecha YYYY-MM-DD a dd/MM/yyyy sin usar Date,
 * para evitar que se muestre el día anterior por zona horaria (UTC vs local).
 */
@Pipe({
  name: 'dateOnly',
  standalone: true
})
export class DateOnlyPipe implements PipeTransform {
  transform(value: string | null | undefined): string {
    if (value == null || value === '') return '';
    const s = String(value).slice(0, 10);
    const parts = s.split('-');
    if (parts.length !== 3) return value;
    const [y, m, d] = parts;
    return `${d}/${m}/${y}`;
  }
}
