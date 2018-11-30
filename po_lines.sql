SELECT DISTINCT hou.name OU,
  pv.vendor_name,
  (case when pra.release_num is null then pha.segment1 else pha.segment1 || '-' || pra.release_num end)  PO_NUM,
  pha.creation_date,
  hla.location_code
FROM apps.po_headers_all pha ,
  apps.po_lines_all pla ,
  apps.po_line_locations_all plla ,
  apps.po_releases_all pra ,
  apps.ap_suppliers pv ,
  apps.hr_organization_units hou ,
  APPS.HR_LOCATIONS_all hla
WHERE pha.vendor_id          = pv.vendor_id
AND pha.po_header_id         = pla.po_header_id
AND pla.po_line_id           = plla.po_line_id
AND pha.org_id               = hou.organization_id
AND plla.ship_to_location_id = hla.location_id
AND plla.PO_RELEASE_ID       = pra.po_release_id(+)
AND pha.creation_date       >= '01-NOV-2018'
AND NVL(pha.CANCEL_FLAG, 'N') = 'N'
ORDER BY pha.creation_date