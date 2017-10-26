package udger

var (
	ClientsSql = `SELECT
			  class_id,
			  client_id,
			  regstring,
			  name,
			  name_code,
			  homepage,
			  icon,
			  icon_big,
			  engine,
			  vendor,
			  vendor_code,
			  vendor_homepage,
			  uptodate_current_version,
			  client_classification,
			  client_classification_code
			FROM udger_client_regex
			  JOIN udger_client_list ON udger_client_list.id = udger_client_regex.client_id
			  JOIN udger_client_class ON udger_client_class.id = udger_client_list.class_id
			ORDER BY sequence
			  ASC`

	OSystemsSql = `SELECT
			  os_id,
			  regstring,
			  family,
			  family_code,
			  name,
			  name_code,
			  homepage,
			  icon,
			  icon_big,
			  vendor,
			  vendor_code,
			  vendor_homepage
			FROM udger_os_regex
			  JOIN udger_os_list ON udger_os_list.id = udger_os_regex.os_id
			ORDER BY sequence
			  ASC`

	DevicesSql = `SELECT
			  deviceclass_id,
			  regstring,
			  name,
			  name,
			  name_code,
			  icon,
			  icon_big
			FROM udger_deviceclass_regex
			  JOIN udger_deviceclass_list ON udger_deviceclass_list.id = udger_deviceclass_regex.deviceclass_id
			ORDER BY sequence
			  ASC`

)
