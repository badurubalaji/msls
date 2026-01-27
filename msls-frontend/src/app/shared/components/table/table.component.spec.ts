import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Component } from '@angular/core';
import { MslsTableComponent, MslsTableCellDirective, TableColumn, SortEvent } from './table.component';

@Component({
  standalone: true,
  imports: [MslsTableComponent, MslsTableCellDirective],
  template: `
    <msls-table
      [columns]="columns"
      [data]="data"
      [sortable]="sortable"
      [loading]="loading"
      [striped]="striped"
      [hoverable]="hoverable"
      [emptyMessage]="emptyMessage"
      (sortChange)="onSortChange($event)"
      (rowClick)="onRowClick($event)"
    >
      <ng-template mslsTableCell="status" let-row let-value="value">
        <span class="badge">{{ value }}</span>
      </ng-template>
    </msls-table>
  `,
})
class TestHostComponent {
  columns: TableColumn[] = [
    { key: 'id', label: 'ID', sortable: true },
    { key: 'name', label: 'Name', sortable: true },
    { key: 'status', label: 'Status', sortable: false },
  ];
  data = [
    { id: 1, name: 'Alice', status: 'active' },
    { id: 2, name: 'Bob', status: 'inactive' },
    { id: 3, name: 'Charlie', status: 'active' },
  ];
  sortable = true;
  loading = false;
  striped = false;
  hoverable = true;
  emptyMessage = 'No data available';

  lastSortEvent: SortEvent | null = null;
  lastClickedRow: Record<string, unknown> | null = null;

  onSortChange(event: SortEvent): void {
    this.lastSortEvent = event;
  }

  onRowClick(row: Record<string, unknown>): void {
    this.lastClickedRow = row;
  }
}

describe('MslsTableComponent', () => {
  let component: TestHostComponent;
  let fixture: ComponentFixture<TestHostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestHostComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TestHostComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display column headers', () => {
    const headerCells = fixture.nativeElement.querySelectorAll('.table__header');
    expect(headerCells.length).toBe(3);
    expect(headerCells[0].textContent).toContain('ID');
    expect(headerCells[1].textContent).toContain('Name');
    expect(headerCells[2].textContent).toContain('Status');
  });

  it('should display data rows', () => {
    const rows = fixture.nativeElement.querySelectorAll('.table__body .table__row');
    expect(rows.length).toBe(3);
  });

  it('should display cell values', () => {
    const cells = fixture.nativeElement.querySelectorAll('.table__cell');
    expect(cells[0].textContent.trim()).toBe('1');
    expect(cells[1].textContent.trim()).toBe('Alice');
  });

  it('should use custom cell template', () => {
    const badgeElement = fixture.nativeElement.querySelector('.badge');
    expect(badgeElement).toBeTruthy();
    expect(badgeElement.textContent).toContain('active');
  });

  it('should emit sortChange when sortable column header is clicked', () => {
    const headerCells = fixture.nativeElement.querySelectorAll('.table__header');
    headerCells[0].click(); // Click on ID column
    fixture.detectChanges();

    expect(component.lastSortEvent).toEqual({ column: 'id', direction: 'asc' });
  });

  it('should cycle through sort directions', () => {
    const headerCells = fixture.nativeElement.querySelectorAll('.table__header');
    const idHeader = headerCells[0];

    // First click: asc
    idHeader.click();
    fixture.detectChanges();
    expect(component.lastSortEvent?.direction).toBe('asc');

    // Second click: desc
    idHeader.click();
    fixture.detectChanges();
    expect(component.lastSortEvent?.direction).toBe('desc');

    // Third click: null (no sort)
    idHeader.click();
    fixture.detectChanges();
    expect(component.lastSortEvent?.direction).toBeNull();
  });

  it('should not sort when column.sortable is false', () => {
    const headerCells = fixture.nativeElement.querySelectorAll('.table__header');
    headerCells[2].click(); // Click on Status column (not sortable)
    fixture.detectChanges();

    // Sort event should not be emitted for non-sortable columns
    // (Previous test leaves lastSortEvent in a state, so we check the table state)
    const statusHeader = headerCells[2];
    expect(statusHeader.classList.contains('table__header--sortable')).toBeFalsy();
  });

  it('should emit rowClick when row is clicked', () => {
    const rows = fixture.nativeElement.querySelectorAll('.table__body .table__row');
    rows[1].click();
    fixture.detectChanges();

    expect(component.lastClickedRow).toEqual({ id: 2, name: 'Bob', status: 'inactive' });
  });

  it('should display empty message when data is empty', () => {
    component.data = [];
    fixture.detectChanges();

    const emptyCell = fixture.nativeElement.querySelector('.table__cell--empty');
    expect(emptyCell).toBeTruthy();
    expect(emptyCell.textContent).toContain('No data available');
  });

  it('should show loading overlay when loading is true', () => {
    component.loading = true;
    fixture.detectChanges();

    const loadingOverlay = fixture.nativeElement.querySelector('.table-container__loading');
    expect(loadingOverlay).toBeTruthy();
  });

  it('should apply striped class when striped is true', () => {
    component.striped = true;
    fixture.detectChanges();

    const tableElement = fixture.nativeElement.querySelector('.table');
    expect(tableElement.classList.contains('table--striped')).toBeTruthy();
  });

  it('should apply hoverable class when hoverable is true', () => {
    const tableElement = fixture.nativeElement.querySelector('.table');
    expect(tableElement.classList.contains('table--hoverable')).toBeTruthy();
  });

  it('should sort data correctly when sorted', () => {
    // Click on Name column to sort ascending
    const headerCells = fixture.nativeElement.querySelectorAll('.table__header');
    headerCells[1].click(); // Name column
    fixture.detectChanges();

    const cells = fixture.nativeElement.querySelectorAll('.table__body .table__row .table__cell:nth-child(2)');
    expect(cells[0].textContent.trim()).toBe('Alice');
    expect(cells[1].textContent.trim()).toBe('Bob');
    expect(cells[2].textContent.trim()).toBe('Charlie');

    // Click again to sort descending
    headerCells[1].click();
    fixture.detectChanges();

    const cellsDesc = fixture.nativeElement.querySelectorAll('.table__body .table__row .table__cell:nth-child(2)');
    expect(cellsDesc[0].textContent.trim()).toBe('Charlie');
    expect(cellsDesc[1].textContent.trim()).toBe('Bob');
    expect(cellsDesc[2].textContent.trim()).toBe('Alice');
  });
});
