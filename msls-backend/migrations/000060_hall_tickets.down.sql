-- Reverse Hall Ticket Generation migration

-- Remove permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'hall-ticket:view', 'hall-ticket:generate', 'hall-ticket:download', 'hall-ticket:template-manage'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'hall-ticket:view', 'hall-ticket:generate', 'hall-ticket:download', 'hall-ticket:template-manage'
);

-- Drop triggers
DROP TRIGGER IF EXISTS set_updated_at_hall_tickets ON hall_tickets;
DROP TRIGGER IF EXISTS set_updated_at_hall_ticket_templates ON hall_ticket_templates;

-- Drop tables (order matters due to references)
DROP TABLE IF EXISTS hall_tickets;
DROP TABLE IF EXISTS hall_ticket_templates;
