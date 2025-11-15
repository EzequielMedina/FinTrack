import { Component, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { SupportDialogComponent } from './support-dialog/support-dialog.component';

interface FAQItem {
  question: string;
  answer: string;
  category: 'general' | 'accounts' | 'cards' | 'transactions' | 'reports';
}

@Component({
  selector: 'app-faq',
  standalone: true,
  imports: [
    CommonModule,
    MatExpansionModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule
  ],
  templateUrl: './faq.component.html',
  styleUrls: ['./faq.component.css']
})
export class FaqComponent {
  private readonly dialog = inject(MatDialog);
  
  selectedCategory = signal<string>('all');

  readonly categories = [
    { value: 'all', label: 'Todas', icon: 'help_outline' },
    { value: 'general', label: 'General', icon: 'info' },
    { value: 'accounts', label: 'Cuentas', icon: 'account_balance' },
    { value: 'cards', label: 'Tarjetas', icon: 'credit_card' },
    { value: 'transactions', label: 'Transacciones', icon: 'swap_horiz' },
    { value: 'reports', label: 'Reportes', icon: 'assessment' }
  ];

  readonly faqs: FAQItem[] = [
    // General
    {
      category: 'general',
      question: '¿Qué es FinTrack?',
      answer: 'FinTrack es una aplicación de gestión financiera personal que te permite administrar tus cuentas, tarjetas, transacciones y generar reportes detallados de tus finanzas.'
    },
    {
      category: 'general',
      question: '¿Es seguro usar FinTrack?',
      answer: 'Sí, FinTrack utiliza encriptación de datos, autenticación segura y cumple con las mejores prácticas de seguridad para proteger tu información financiera.'
    },
    {
      category: 'general',
      question: '¿Puedo usar FinTrack en múltiples dispositivos?',
      answer: 'Sí, puedes acceder a FinTrack desde cualquier dispositivo con un navegador web. Tus datos se sincronizan automáticamente.'
    },
    
    // Cuentas
    {
      category: 'accounts',
      question: '¿Qué tipos de cuentas puedo crear?',
      answer: 'Puedes crear Billeteras Virtuales (para pagos digitales), Cuentas Bancarias (que pueden tener múltiples tarjetas), Cuentas de Ahorro y Cuentas Corrientes.'
    },
    {
      category: 'accounts',
      question: '¿Cómo agrego una nueva cuenta?',
      answer: 'Ve a la sección "Cuentas", haz clic en "Nueva Cuenta", selecciona el tipo de cuenta, completa los datos requeridos y confirma. La cuenta se creará inmediatamente.'
    },
    {
      category: 'accounts',
      question: '¿Puedo tener cuentas en diferentes monedas?',
      answer: 'Sí, actualmente soportamos cuentas en Pesos Argentinos (ARS) y Dólares Estadounidenses (USD). Puedes crear múltiples cuentas en diferentes monedas.'
    },
    {
      category: 'accounts',
      question: '¿Cómo se calcula el límite de crédito total?',
      answer: 'El límite de crédito total es la suma de todos los límites de tus tarjetas de crédito activas, incluyendo las tarjetas asociadas a tus cuentas bancarias.'
    },
    
    // Tarjetas
    {
      category: 'cards',
      question: '¿Cómo agrego una tarjeta?',
      answer: 'Ve a "Tarjetas", haz clic en "Agregar Tarjeta", selecciona la cuenta asociada, completa los datos de la tarjeta (número, titular, vencimiento) y confirma. La información se almacena de forma segura.'
    },
    {
      category: 'cards',
      question: '¿Puedo tener múltiples tarjetas en una cuenta?',
      answer: 'Sí, las Cuentas Bancarias y Cuentas Corrientes pueden tener múltiples tarjetas (tanto de crédito como de débito) asociadas.'
    },
    {
      category: 'cards',
      question: '¿Qué hago si pierdo o me roban una tarjeta?',
      answer: 'Ve a la sección "Tarjetas", encuentra la tarjeta comprometida, haz clic en el menú de opciones y selecciona "Bloquear". Esto desactivará la tarjeta inmediatamente.'
    },
    {
      category: 'cards',
      question: '¿Puedo establecer una tarjeta como predeterminada?',
      answer: 'Sí, puedes marcar una tarjeta como predeterminada. Esta será la tarjeta sugerida automáticamente al realizar transacciones desde esa cuenta.'
    },
    
    // Transacciones
    {
      category: 'transactions',
      question: '¿Cómo registro una transacción?',
      answer: 'Ve a "Transacciones", haz clic en "Nueva Transacción", selecciona el tipo (ingreso, gasto, transferencia), completa los detalles y confirma. La transacción se registrará y actualizará el balance automáticamente.'
    },
    {
      category: 'transactions',
      question: '¿Puedo editar o eliminar transacciones?',
      answer: 'Sí, puedes editar o eliminar transacciones que hayas registrado. Ve a la lista de transacciones, haz clic en la transacción deseada y selecciona "Editar" o "Eliminar".'
    },
    {
      category: 'transactions',
      question: '¿Qué tipos de transacciones puedo registrar?',
      answer: 'Puedes registrar Ingresos (dinero que recibes), Gastos (dinero que gastas), Transferencias (entre tus cuentas), Pagos con tarjeta y Pagos de cuotas.'
    },
    {
      category: 'transactions',
      question: '¿Las transacciones se pueden categorizar?',
      answer: 'Sí, puedes asignar categorías a tus transacciones (alimentación, transporte, entretenimiento, etc.) para un mejor seguimiento y análisis en los reportes.'
    },
    
    // Reportes
    {
      category: 'reports',
      question: '¿Qué tipos de reportes puedo generar?',
      answer: 'Puedes generar reportes de Transacciones (detallado por período), Gastos vs Ingresos (comparativo mensual), Gastos por Categoría y Saldos de Cuentas.'
    },
    {
      category: 'reports',
      question: '¿Puedo exportar los reportes?',
      answer: 'Sí, todos los reportes se pueden exportar en formato PDF. Haz clic en el botón "Exportar PDF" en la esquina superior derecha del reporte.'
    },
    {
      category: 'reports',
      question: '¿Los reportes se actualizan en tiempo real?',
      answer: 'Sí, los reportes se generan con los datos más recientes de tu cuenta. Cada vez que generas un reporte, obtienes la información actualizada.'
    },
    {
      category: 'reports',
      question: '¿Puedo filtrar reportes por período?',
      answer: 'Sí, puedes filtrar todos los reportes por rango de fechas personalizado, mes actual, mes anterior, últimos 3 meses, últimos 6 meses o año completo.'
    }
  ];

  get filteredFaqs(): FAQItem[] {
    const category = this.selectedCategory();
    if (category === 'all') {
      return this.faqs;
    }
    return this.faqs.filter(faq => faq.category === category);
  }

  selectCategory(category: string): void {
    this.selectedCategory.set(category);
  }

  openSupportDialog(): void {
    const dialogRef = this.dialog.open(SupportDialogComponent, {
      width: '600px',
      maxWidth: '95vw',
      disableClose: false,
      autoFocus: true
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log('Email de soporte enviado exitosamente');
      }
    });
  }
}
