===========================
Streaming data from the CLI
===========================

The following example creates a budget for a resource named ``storage-account`` and a
rule that will trigger if 135% or more of the monthly budget amount has been reached. 

.. code-block:: bash

  # create budget / rules
  cat <<EOF > config.yaml
  budgets:
  - resource_id: storage-account
    monthly_amount: 100
    rules:
    - type: percentage
      name: early-warning
      value: 135

  EOF
  
  # execute
  wrangler -c ./config.yaml

When executed, the application will prompt for input:

.. code-block:: bash

  2024-02-29 18:38:14 INF Reading data from stdin...

You can then enter CSV data from the CLI. To demonstrate, enter the following:

.. code-block:: bash

  2024-02-05T12:00:00Z,storage-account,monthly,110
  2024-02-05T12:00:00Z,storage-account,daily,14
  2024-02-06T12:00:00Z,storage-account,monthly,150
  2024-02-06T12:00:00Z,storage-account,daily,21

Then enter `Ctrl+D` to indicate the end of input. Wrangler will evaluate each input
against the budget rules we created earlier and it should fail with a single violation:

.. code-block:: bash

  2024-02-29 18:38:37 INF Evaluating 4 record(s)...
  ...
  2024-02-29 18:38:37 INF 1 violation(s) found
  +-----------------+---------------+--------------------------------+------------+---------------+---------------+
  | Resource ID     | Rule name     | Condition                      | Date       | Budget amount | Actual amount |
  +-----------------+---------------+--------------------------------+------------+---------------+---------------+
  | storage-account | early-warning | actual amount >= 135.00% of    | 2024-03-06 | 100.00        | 150.00        |
  |                 |               | budget                         |            |               |               |
  +-----------------+---------------+--------------------------------+------------+---------------+---------------+

  2024-02-29 18:38:37 FTL Failed with violations

