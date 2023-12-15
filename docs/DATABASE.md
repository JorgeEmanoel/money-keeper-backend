users
    id
    name
    email
    password

plans
    id
    name
    description
    status (ENABLED, DISABLED)
    user_id

skeletons
    id
    name
    description
    direction (income, outcome)
    frequency (monthly, anual, random)
    value
    currency
    plan_id
    user_id

transactions
    id
    name
    description
    direction (income, outcome)
    value
    currency
    reference (MONTH-YEAR)
    status (PENDING, PAID, CANCELED)
    user_id
