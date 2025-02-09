export function useValidation() {
  function validate(fields) {
    let errors = {};
    for (let key in fields) {
      // Check if the value is null, undefined, or an empty string (but allow 0)
      if (
        fields[key] === null ||
        fields[key] === undefined ||
        (typeof fields[key] === "string" && fields[key].trim() === "")
      ) {
        // Remove '_id' suffix if it exists and replace underscores with spaces
        let fieldName = key.replace(/_id$/, "").replace(/_/g, " ");
        errors[key] = `The ${fieldName} field is required.`;
      }
    }
    return errors;
  }

  function validatePasswordMatch(password, confirm_password) {
    let error = {};
    if (password !== confirm_password) {
      error["confirm_password"] =
        "The passwords you entered don’t match. Please try again.";
    }
    return error;
  }

  return { validate, validatePasswordMatch };
}
