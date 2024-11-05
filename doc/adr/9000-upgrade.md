# 9000. Upgrade skupper sites from v1 to v2

Date: 2024-11-04

## Status

Proposed

## Context

We need to document how customers will upgrade v1 skupper sites to
v2 skupper sites.

(TODO: refer to some other doc that explains: v2 configuration via CRs.)

## High level decisions

High level decisions to document further:

1. The skupper community will define the recommended procedure
   to upgrade v1 skupper sites to v2 skupper sites.

2. skupper v1 and skupper v2 cannot be installed in the same
   namespace at the same time.

   Reasons include:
     * TODO: verify: configmapp name collision?
     * TODO: verify: skupper service collision?
     * TODO: verify: secret name collision?

3. A tool will be created to create skupper v2 CRs for given
   skupper v1 sites. 

   Details include:
     *  

4. 

5. 

## Open upgrade questions

1. aaa 

2. bbb

3. ccc

## Fixed vs immutable resources during upgrade

Fixed resources include:

1. Namespace.  

Mutable resources include:

2. Annotated services.

## Options

## Decision

## Consequences

## HCM open questions

HCM open questions include:

## TODOs
