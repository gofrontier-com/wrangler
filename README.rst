.. image:: https://pkg.go.dev/badge/github.com/gofrontier-com/wrangler.svg
    :target: https://pkg.go.dev/github.com/gofrontier-com/wrangler
.. image:: https://github.com/gofrontier-com/wrangler/actions/workflows/ci.yml/badge.svg
    :target: https://github.com/gofrontier-com/wrangler/actions/workflows/ci.yml

========
Wrangler
========

Wrangler is a command line tool for cost management.

.. contents:: Table of Contents
    :local:

-----
About
-----

Wrangler was built to simplify cloud cost management. In the world of distributed cloud applications, 
tracking costs can be challenging. With Wrangler, teams gain control and visibility over costs by
centralising budget configuration, streamling trigger rules and executing configurable actions when 
thresholds are reached.

Consuming data
##############

By default Wrangler will attempt to read data from the CLI standard input in CSV format::

    timestamp, resource_id, period, value, baseline, currency, category

- **timestamp** - Date/time of charge (Unix timestamp or ISO 8601 date/time)
- **resource_id** - Unique name/ID of the resource that accrued the charge (any)
- **period** - Period the charge covers (one of [``daily``, ``monthly``])
- **value** - Amount accrued for period (decimal)
- **baseline** - Optional, if specified then allows for dynamic budgeting (decimal)
- **currency** - Optional, if specified then only budgets or rules of qualifying currencies will apply (ISO-4217 code)
- **category** - Optional, if specified can be used to group budgets / rules (any)

Data can also be consumed through a `custom provider <#providers>`_. If you want to disable the default
behaviour use the ``--no-stdin`` flag.

Some examples of the types of use-cases Wrangler can be used for are:

- Flagging over and underspend of resources
- Forecasting future overspend
- Detecting cost anomalies (e.g. unusual spike or reduction in cost over a defined period of time)
- Automating budgets by analysing previous spend

.. _providers:

Custom providers
################

TODO

-------------
Configuration
-------------

Examples
--------

Monthly budget

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100
    ...

Daily budget

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      daily_amount: 10
    ...


Monthly & daily budget

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100
      daily_amount: 10
    ...

Currency-specific budgets

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100
      currency: GBP
    - resource_id: my-resource
      monthly_amount: 140
      currency: USD
    ...

*note - rules will only trigger for cost data matching same currency*

Budget with local rule

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100
      rules:
      - name: nearing-budget-alert
        type: percentage
        value: 85
    ...

Budget with global rule

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100

    rules:
    - name: overspend-alert
      type: percentage
      value: 135
    ...

Budget with global rule filter

.. code:: yaml
    ---
    budgets:
    - resource_id: my-resource
      monthly_amount: 100

    rules:
    - name: infra-overspend-alert
      type: percentage
      value: 135
      categories: 
      - infra
    ...

*note - rule will only trigger for cost data matching same categories*

Scenarios
---------

TODO

.. _trigger-rules:

Trigger rules
-------------

Fixed amount
------------

Percentage
----------

Overrun
-------

--------
Download
--------

Binaries and packages of the latest stable release are available at `https://github.com/gofrontier-com/wrangler/releases <https://github.com/gofrontier-com/wrangler/releases>`_.

-----
Usage
-----

.. code:: bash

  $ wrangler --help
    WWrangler is a command line tool for cost management.

    Usage:
      wrangler [flags]

    Flags:
      -c, --config string   configuration file (default "<cwd>/.wrangler.yaml")
      -h, --help            help for wrangler
      -v, --version         version for wrangler

Stream data from standard input:

.. code:: bash

    $ echo "..." | wrangler

Stream data from a file:

.. code:: bash

    $ cat costdata.csv | wrangler

Disable streaming:

.. code:: bash

    $ wrangler --no-stdin

------------
Contributing
------------

We welcome contributions to this repository. Please see `CONTRIBUTING.md <https://github.com/gofrontier-com/wrangler/tree/main/CONTRIBUTING.md>`_ for more information.
