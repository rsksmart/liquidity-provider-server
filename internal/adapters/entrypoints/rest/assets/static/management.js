<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Management UI</title>
    <link href="./static/Bootstrap.css" rel="stylesheet" crossorigin="anonymous">
    <link href="./static/management.css" rel="stylesheet" />
</head>
<body>
    <div class="container main-content">
        <div class="row">
            <div class="col-md-12">
                <h1>Management Dashboard</h1><hr>
            </div>
        </div>
        <div class="row compact-row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">Provider</div>
                    <div class="card-body">
                        <h5 class="card-title">Provider RSK Address</h5>
                        <p class="card-text" id="providerRskAddress"></p>
                        <h5 class="card-title">Provider BTC Address</h5>
                        <p class="card-text" id="providerBtcAddress"></p>
                        <h5 class="card-title">Operational Status</h5>
                        <p class="card-text" id="isOperational"></p>
                    </div>
                </div>
                <div class="card">
                    <div class="card-header">Collateral</div>
                    <div class="card-body">
                        <ul class="nav nav-tabs" id="collateralTabs" role="tablist">
                            <li class="nav-item" role="presentation">
                                <a class="nav-link active" id="pegin-tab" data-bs-toggle="tab" href="#pegin" role="tab" aria-controls="pegin" aria-selected="true">Pegin</a>
                            </li>
                            <li class="nav-item" role="presentation">
                                <a class="nav-link" id="pegout-tab" data-bs-toggle="tab" href="#pegout" role="tab" aria-controls="pegout" aria-selected="false">Pegout</a>
                            </li>
                        </ul>
                        <div class="tab-content" id="collateralTabContent">
                            <div class="tab-pane fade show active" id="pegin" role="tabpanel" aria-labelledby="pegin-tab">
                                <h5 class="card-title">Pegin Collateral</h5>
                                <p class="card-text" id="peginCollateral"></p>
                                <div class="collateral-inputs">
                                    <div class="mb-3">
                                        <label for="addPeginCollateralAmount" class="form-label">Add Pegin Collateral Amount</label>
                                        <input type="number" class="form-control" id="addPeginCollateralAmount" placeholder="Enter amount in rBTC">
                                    </div>
                                </div>
                                <div class="collateral-buttons">
                                    <button type="button" class="btn btn-primary" id="addPeginCollateralButton">Add Pegin Collateral</button>
                                    <div class="loading-bar" id="peginLoadingBar"></div>
                                </div>
                            </div>
                            <div class="tab-pane fade" id="pegout" role="tabpanel" aria-labelledby="pegout-tab">
                                <h5 class="card-title">Pegout Collateral</h5>
                                <p class="card-text" id="pegoutCollateral"></p>
                                <div class="collateral-inputs">
                                    <div class="mb-3">
                                        <label for="addPegoutCollateralAmount" class="form-label">Add Pegout Collateral Amount</label>
                                        <input type="number" class="form-control" id="addPegoutCollateralAmount" placeholder="Enter amount in rBTC">
                                    </div>
                                </div>
                                <div class="collateral-buttons">
                                    <button type="button" class="btn btn-primary" id="addPegoutCollateralButton">Add Pegout Collateral</button>
                                    <div class="loading-bar" id="pegoutLoadingBar"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">Configuration</div>
                    <div class="card-body">
                        <h5 class="card-title">Current Configuration</h5>
                        <ul class="nav nav-tabs" id="configTabs" role="tablist">
                            <li class="nav-item" role="presentation">
                                <a class="nav-link active" id="general-tab" data-bs-toggle="tab" href="#general" role="tab" aria-controls="general" aria-selected="true">General</a>
                            </li>
                            <li class="nav-item" role="presentation">
                                <a class="nav-link" id="peginConfig-tab" data-bs-toggle="tab" href="#peginConfig" role="tab" aria-controls="peginConfig" aria-selected="false">Pegin</a>
                            </li>
                            <li class="nav-item" role="presentation">
                                <a class="nav-link" id="pegoutConfig-tab" data-bs-toggle="tab" href="#pegoutConfig" role="tab" aria-controls="pegoutConfig" aria-selected="false">Pegout</a>
                            </li>
                        </ul>
                        <div class="tab-content" id="configTabContent">
                            <div class="tab-pane fade show active" id="general" role="tabpanel" aria-labelledby="general-tab">
                                <div id="generalConfig"></div>
                            </div>
                            <div class="tab-pane fade" id="peginConfig" role="tabpanel" aria-labelledby="peginConfig-tab">
                                <div id="peginConfig"></div>
                            </div>
                            <div class="tab-pane fade" id="pegoutConfig" role="tabpanel" aria-labelledby="pegoutConfig-tab">
                                <div id="pegoutConfig"></div>
                            </div>
                        </div>
                        <button type="button" class="btn btn-primary" id="saveConfig">Save Configuration</button>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="toast-container">
        <div id="successToast" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-header">
                <strong class="me-auto">Success</strong>
                <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
            <div class="toast-body">
                Configuration saved successfully!
            </div>
        </div>
    </div>

    <script src="./static/Bootstrap.js" crossorigin="anonymous"></script>
    <script src="./static/decimal.js" crossorigin="anonymous"></script>
    <script nonce="{{ .ScriptNonce }}">
        const data = {{.}};
    </script>
    <script type="module" src="./static/configUtils.js" crossorigin="anonymous"></script>
    <script type="module" src="./static/management.js" crossorigin="anonymous"></script>
</body>
</html>