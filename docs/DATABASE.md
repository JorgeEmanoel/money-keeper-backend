users
    id
    name
    email
    password
    main_income_value
    main_income_currency

transactions
    id
    name
    description
    type (income, outcome)
    frequency (monthly, anual, random)
    value
    currency
    reference (MONTH-YEAR)
    user_id

commitments
    id
    transaction_id
    reference (MONTH-YEAR)
