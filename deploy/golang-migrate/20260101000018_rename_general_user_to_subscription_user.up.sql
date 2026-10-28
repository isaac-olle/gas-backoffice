UPDATE backoffice.users
SET system_role = 'subscription_user'
WHERE system_role = 'general_user';

