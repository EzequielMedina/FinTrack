import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Pipe({
  name: 'markdown',
  standalone: true
})
export class MarkdownPipe implements PipeTransform {
  constructor(private sanitizer: DomSanitizer) {}

  transform(value: string): SafeHtml {
    if (!value) return '';

    let html = value;

    // Headers (### Título)
    html = html.replace(/^### (.*$)/gim, '<h3 class="md-h3">$1</h3>');
    html = html.replace(/^## (.*$)/gim, '<h2 class="md-h2">$1</h2>');
    html = html.replace(/^# (.*$)/gim, '<h1 class="md-h1">$1</h1>');

    // Bold (**texto**)
    html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="md-bold">$1</strong>');
    
    // Italic (*texto*)
    html = html.replace(/\*(.*?)\*/g, '<em class="md-italic">$1</em>');

    // Montos ($X,XXX.XX)
    html = html.replace(/\$([0-9,]+\.?[0-9]*)/g, '<span class="md-amount">$$$1</span>');

    // Listas con guiones (- Item)
    const lines = html.split('\n');
    let inList = false;
    const processedLines: string[] = [];

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const listMatch = line.match(/^[\s]*[-•]\s+(.+)$/);
      
      if (listMatch) {
        if (!inList) {
          processedLines.push('<ul class="md-list">');
          inList = true;
        }
        processedLines.push(`<li class="md-list-item">${listMatch[1]}</li>`);
      } else {
        if (inList) {
          processedLines.push('</ul>');
          inList = false;
        }
        processedLines.push(line);
      }
    }
    
    if (inList) {
      processedLines.push('</ul>');
    }

    html = processedLines.join('\n');

    // Line breaks - reducir saltos de línea múltiples a uno solo
    html = html.replace(/\n{2,}/g, '\n'); // Máximo 1 salto (eliminar dobles)
    html = html.replace(/\n/g, '<br>');

    // Links [texto](url)
    html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" class="md-link">$1</a>');

    // Emojis destacados (para que no pierdan su espacio)
    html = html.replace(/(💰|💳|📊|💸|📈|📉|✅|❌|⚠️|🔥|💡|📅)/g, '<span class="md-emoji">$1</span>');

    return this.sanitizer.bypassSecurityTrustHtml(html);
  }
}
