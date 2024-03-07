package pegin

const (
	// CallForUserExtraGas
	/**
	 *	The gas spent in the callForUser function can be divided in two parts,
	 *	the first part is the gas spent in the callForUser function itself, call done on behalf
	 *	of the user. This constant represents the first part and needs to be added to the estimation
	 *	done during the get pegin quote process.
	 */
	CallForUserExtraGas = 180000
)
