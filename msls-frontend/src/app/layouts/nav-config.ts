import { IconName } from '../shared/components/icon/icon.component';

/**
 * Navigation item interface for sidebar navigation.
 */
export interface NavItem {
  /** Display label for the navigation item */
  label: string;

  /** Icon name (Heroicons outline) */
  icon: IconName;

  /** Route path for navigation (optional if has children) */
  route?: string;

  /** Child navigation items for nested menus */
  children?: NavItem[];

  /** Required permissions to see this item */
  permissions?: string[];

  /** Whether this item is disabled */
  disabled?: boolean;

  /** Badge text to display (e.g., "New", "Beta") */
  badge?: string;
}

/**
 * Main navigation items for the application sidebar.
 */
export const NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    icon: 'home',
    route: '/dashboard',
  },
  {
    label: 'Students',
    icon: 'users',
    route: '/students',
  },
  {
    label: 'Staff',
    icon: 'briefcase',
    children: [
      {
        label: 'Directory',
        icon: 'users',
        route: '/staff',
      },
      {
        label: 'Documents',
        icon: 'document-text',
        route: '/staff-documents/compliance',
        permissions: ['staff_document.view'],
      },
      {
        label: 'Expiring Docs',
        icon: 'exclamation-triangle',
        route: '/staff-documents/expiring',
        permissions: ['staff_document.view'],
      },
      {
        label: 'Document Types',
        icon: 'cog',
        route: '/staff-documents/types',
        permissions: ['staff_document.manage_types'],
      },
    ],
  },
  {
    label: 'Payroll',
    icon: 'banknotes',
    route: '/payroll',
    permissions: ['payroll.view'],
  },
  {
    label: 'Attendance',
    icon: 'clock',
    permissions: ['attendance:mark_self', 'attendance:view_all', 'student_attendance:mark_class'],
    children: [
      {
        label: 'My Attendance',
        icon: 'user',
        route: '/attendance/my',
        permissions: ['attendance:mark_self'],
      },
      {
        label: 'Manage Staff',
        icon: 'clipboard-document-list',
        route: '/attendance/manage',
        permissions: ['attendance:view_all'],
      },
      {
        label: 'Mark Student',
        icon: 'user-group',
        route: '/student-attendance/mark',
        permissions: ['student_attendance:mark_class'],
      },
      {
        label: 'Period-wise',
        icon: 'table-cells',
        route: '/student-attendance/period',
        permissions: ['student_attendance:mark_class'],
      },
      {
        label: 'Daily Summary',
        icon: 'chart-bar',
        route: '/student-attendance/summary',
        permissions: ['student_attendance:view_class'],
      },
    ],
  },
  {
    label: 'Academics',
    icon: 'academic-cap',
    children: [
      {
        label: 'Classes',
        icon: 'user-group',
        route: '/academics/classes',
      },
      {
        label: 'Sections',
        icon: 'rectangle-stack',
        route: '/academics/sections',
      },
      {
        label: 'Streams',
        icon: 'arrows-right-left',
        route: '/academics/streams',
      },
      {
        label: 'Subjects',
        icon: 'book-open',
        route: '/academics/subjects',
      },
      {
        label: 'Assignments',
        icon: 'clipboard-document-list',
        route: '/assignments',
        permissions: ['assignment.view'],
      },
      {
        label: 'Workload',
        icon: 'chart-bar',
        route: '/assignments/workload',
        permissions: ['assignment.workload'],
      },
      {
        label: 'Timetable',
        icon: 'calendar',
        route: '/academics/timetable',
        children: [
          {
            label: 'Overview',
            icon: 'calendar',
            route: '/academics/timetable',
          },
          {
            label: 'Timetables',
            icon: 'table-cells',
            route: '/academics/timetable/list',
          },
          {
            label: 'My Timetable',
            icon: 'user',
            route: '/academics/timetable/my',
          },
          {
            label: 'Shifts',
            icon: 'clock',
            route: '/academics/timetable/shifts',
          },
          {
            label: 'Day Patterns',
            icon: 'calendar-days',
            route: '/academics/timetable/day-patterns',
          },
          {
            label: 'Period Slots',
            icon: 'table-cells',
            route: '/academics/timetable/period-slots',
          },
          {
            label: 'Substitutions',
            icon: 'arrows-right-left',
            route: '/academics/timetable/substitutions',
            permissions: ['substitution:view'],
          },
        ],
      },
    ],
  },
  {
    label: 'Exams',
    icon: 'file-signature',
    permissions: ['exam:type:view'],
    children: [
      {
        label: 'Exam Types',
        icon: 'clipboard-document-list',
        route: '/exams/types',
        permissions: ['exam:type:view'],
      },
    ],
  },
  {
    label: 'Admissions',
    icon: 'clipboard-document-list',
    permissions: ['admissions:read', 'enquiries:read'],
    children: [
      {
        label: 'Overview',
        icon: 'information-circle',
        route: '/admissions',
        permissions: ['admissions:read'],
      },
      {
        label: 'Dashboard',
        icon: 'chart-bar',
        route: '/admissions/dashboard',
        permissions: ['admissions:read'],
      },
      {
        label: 'Sessions',
        icon: 'calendar-days',
        route: '/admissions/sessions',
        permissions: ['admissions:read'],
      },
      {
        label: 'Enquiries',
        icon: 'comment-dots',
        route: '/admissions/enquiries',
        permissions: ['enquiries:read'],
      },
      {
        label: 'Applications',
        icon: 'file-pen',
        route: '/admissions/applications',
        permissions: ['admissions:read'],
      },
      {
        label: 'Tests',
        icon: 'list-check',
        route: '/admissions/tests',
        permissions: ['admissions:read'],
      },
      {
        label: 'Merit List',
        icon: 'trophy',
        route: '/admissions/merit-list',
        permissions: ['admissions:read'],
      },
    ],
  },
  {
    label: 'Finance',
    icon: 'currency-dollar',
    children: [
      {
        label: 'Fee Collection',
        icon: 'banknotes',
        route: '/finance/fees',
      },
      {
        label: 'Reports',
        icon: 'chart-bar',
        route: '/finance/reports',
      },
    ],
  },
  {
    label: 'Settings',
    icon: 'cog',
    route: '/settings',
  },
  {
    label: 'Admin',
    icon: 'shield-check',
    permissions: ['roles:read', 'settings:write', 'branches:read', 'academic-years:read', 'department:read', 'designation:read'],
    children: [
      {
        label: 'Branches',
        icon: 'map-pin',
        route: '/admin/branches',
        permissions: ['branches:read'],
      },
      {
        label: 'Academic Years',
        icon: 'calendar',
        route: '/admin/academic-years',
        permissions: ['academic-years:read'],
      },
      {
        label: 'Departments',
        icon: 'building-office',
        route: '/admin/departments',
        permissions: ['department:read'],
      },
      {
        label: 'Designations',
        icon: 'identification',
        route: '/admin/designations',
        permissions: ['designation:read'],
      },
      {
        label: 'Roles',
        icon: 'user-group',
        route: '/admin/roles',
        permissions: ['roles:read'],
      },
      {
        label: 'Feature Flags',
        icon: 'flag',
        route: '/admin/feature-flags',
        permissions: ['settings:write'],
      },
    ],
  },
];
