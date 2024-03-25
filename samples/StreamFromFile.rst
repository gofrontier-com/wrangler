==========================
Streaming data from a file
==========================

The following example creates a budget for a resource named ``storage-account`` and a
rule that will trigger if 100% or more of the daily budget amount has been reached. 

The sample data has two records that violate this rule and therefore should trigger when executed.

.. code-block:: bash

  # create budget / rules
  cat <<EOF > config.yaml
  budgets:
  - resource_id: storage-account
    monthly_amount: 1000
    daily_amount: 20
    rules:
    - name: budget-reached
      type: percentage
      period: daily
      value: 100

  EOF

  # create sample data
  echo "1707044241,storage-account,daily,14" > data.csv
  echo "1707138019,storage-account,daily,5.5" >> data.csv
  echo "1707376221,storage-account,daily,50" >> data.csv
  echo "1707743640,storage-account,daily,20" >> data.csv
  echo "1707889825,storage-account,monthly,1050" >> data.csv
  
  # execute
  cat ./data.csv | wrangler -c ./config.yaml

When run, the application should fail with 2 violations:

.. code-block:: bash

  2024-02-29 18:01:18 INF 2 violation(s) found
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+
  | Resource ID     | Rule name      | Condition                      | Date       | Budget amount | Actual amount |
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+
  | storage-account | budget-reached | actual amount >= 100.00% of    | 2024-02-08 | 20.00         | 50.00         |
  |                 |                | budget                         |            |               |               |
  | storage-account | budget-reached | actual amount >= 100.00% of    | 2024-02-12 | 20.00         | 20.00         |
  |                 |                | budget                         |            |               |               |
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+
  
  2024-02-29 18:01:19 FTL Failed with violations


Lets adjust our rule to handle both daily & monthly budgets. Open up ``config.yaml`` and delete the ``period`` field
from the rule. Your config should now look like:

.. code-block:: bash

  budgets:
  - resource_id: storage-account
    monthly_amount: 1000
    daily_amount: 20
    rules:
    - name: budget-reached
      type: percentage
      value: 100

Save the config and re-run wrangler, it should now fail with 3 violations:

.. code-block:: bash

  2024-02-29 18:01:18 INF 3 violation(s) found
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+
  | Resource ID     | Rule name      | Condition                      | Date       | Budget amount | Actual amount |
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+
  | storage-account | budget-reached | actual amount >= 100.00% of    | 2024-02-08 | 20.00         | 50.00         |
  |                 |                | budget                         |            |               |               |
  | storage-account | budget-reached | actual amount >= 100.00% of    | 2024-02-12 | 20.00         | 20.00         |
  |                 |                | budget                         |            |               |               |
  | storage-account | budget-reached | actual amount >= 100.00% of    | 2024-02-14 | 1000.00       | 1050.00       |
  |                 |                | budget                         |            |               |               |
  +-----------------+----------------+--------------------------------+------------+---------------+---------------+

  2024-02-29 18:01:19 FTL Failed with violations