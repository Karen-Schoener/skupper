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

2. A tool will be created to create skupper v2 CRs for given
   skupper v1 sites. 

   Details include:
     * In order to read the skupper v1 configuration, the tool
       will be provided online access to the skupper v1 cluster.

3. skupper v1 and skupper v2 cannot be installed in the same
   namespace at the same time.

   Reasons include:
     * TODO: verify: configmap name collision?
     * TODO: verify: service name collision?
     * TODO: verify: secret name collision?

4. v2 link creation

   To create links, the user will:
     * create AccessGrant CRs in the target v2 site
     * wait for the v2 Controller to populate the
       AccessGrant CR with: CA, code, URL   
     * edit an AccessToken CR with: CA, code, URL
     * create AccessToken CRs in the originating v2 site

   Discuss: 
     * Are links transferable?
     * Is there a low friction way to create v2 links for the user?

5. Upgrade steps: delete v1 site, create v2 site in same namespace.

### Upgrade steps: delete v1 site, create v2 site in same namespace

Upgrade steps include:

  1. Provide online access to all v1 sites to tool.

     * Tool reads configmap skupper-site for all sites.

       Tool keeps a map of site uid and site config.

     * Tool identifies configured links by: reading secret
       metadata type=connection-token.

     * Tool identifies service configuration by: reading configmap
       skupper-servivces.

  2. Tool generates CRs.

     CRs may be stored as yaml files.

     If upgrade tool is called programatically, CRs may be returned
     programatically as yaml strings.

     * Site CR

       For every v1 skupper-site configmap

     * AccessGrant CR

       For every v1 secret type=connection-token.

       Use generated-by annotation to identify target site.

       TODO: discuss: 1 access grant per v1 link?

     * AccessToken CR

       For every v1 secret type=connection-token.

       Note: credential fields are populated at a later step.

     * Listener CR

       For every skupper-serices configmap entry.

       Create a Listener CR for each service port.

     * Connector CR

       For every skupper-serices configmap entry with targets populated.

       Create a Connector CR for each target entry and for each target port.

  3. User deletes v1 skupper sites.

  4. User installs v2 CRDs in cluster.
     This step requires cluster permissions.

  5. User creates v2 skupper controller deployment.
     User starts v2 skupper controller.

     * User may start 1 controller per cluster.
     * User may start 1 controller per namespace.

  6. User applies v2 yaml files:

     * Site CR
     * Access Grant CRs
     * Listener CRs
     * Connector CRs

  7. User waits for v2 Controller to populate AccessGrant CRs with
     credential info: CA, code, URL.

  8. User edits v2 AccessToken CRs to include AccessGrant credential
     info: CA, code, URL.
  
## Open upgrade questions

1. How will tool read v1 skupper state: online (kube access) vs 
   offline (debug dump).

2. Identify upgrade options:

   * delete v1 site, create v2 site in same namespace.

   * create v2 site in same namespace.

     v1 and v2 run in the same namespace for some time period.

   * create v2 site in separate namespace.  

     v1 and v2 sites both run in cluster for some time period.

3. 

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
